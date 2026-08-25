package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	//"fmt"

	"MeshNet/internal/domain"

	_ "modernc.org/sqlite"
)

type record struct {
	Object  *domain.Object  `json:"object"`
	Payload *domain.Payload `json:"payload"`
}

type Store struct {
	db *sql.DB
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	store := &Store{db: db}

	if err := store.init(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) init() error {
	_, err := s.db.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS objects (
			id		TEXT PRIMARY KEY,
			hash	TEXT NOT NULL UNIQUE,
			object 	BLOB NOT NULL,
			payload	BLOB
		);
		`)
	return err
}

func isUniqueViolation(err error) bool {
	var sqliteErr interface {
		ErrorCode() int
	}

	if !errors.As(err, &sqliteErr) {
		return false
	}

	const sqliteConstraint = 19

	return sqliteErr.ErrorCode() == sqliteConstraint
}

func (s *Store) Save(ctx context.Context, obj *domain.Object, payload *domain.Payload) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if obj == nil {
		return domain.ErrNotFound
	}

	objectData, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	var payloadData []byte

	if payload != nil {
		payloadData, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO objects (
			id,
			hash,
			object,
			payload
		)
		VALUES (?, ?, ?, ?)
	`, obj.ID, obj.Hash, objectData, payloadData)

	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrDuplicate
		}

		return err
	}
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	var objectData 	[]byte
	var payloadData []byte

	err := s.db.QueryRowContext(ctx, `
		SELECT object, payload
		FROM objects
		WHERE id = ?
	`, id).Scan(&objectData, &payloadData)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, domain.ErrNotFound
		}
		return nil, nil, err
	}

	var obj *domain.Object
	if err := json.Unmarshal(objectData, &obj); err != nil {
		return nil, nil, err
	}

	var payload *domain.Payload
	if len(payloadData) > 0 {
		payload = new(domain.Payload)

		if err := json.Unmarshal(payloadData, payload); err != nil {
			return nil, nil, err
		}
	}

	return obj, payload, nil
}

func (s *Store) GetByHash(ctx context.Context, hash string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	var objectData 	[]byte
	var payloadData []byte

	err := s.db.QueryRowContext(ctx, `
		SELECT object, payload
		FROM objects
		WHERE hash = ?
	`, hash).Scan(&objectData, &payloadData)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, domain.ErrNotFound
		}
		return nil, nil, err
	}

	var obj *domain.Object
	if err := json.Unmarshal(objectData, &obj); err != nil {
		return nil, nil, err
	}

	var payload *domain.Payload
	if len(payloadData) > 0 {
		payload = new(domain.Payload)

		if err := json.Unmarshal(payloadData, payload); err != nil {
			return nil, nil, err
		}
	}

	return obj, payload, nil
}

func (s *Store) List(ctx context.Context) ([]*domain.Object, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT object
		FROM objects
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := make([]*domain.Object, 0)

	for rows.Next() {
		var objectData []byte

		if err := rows.Scan(&objectData); err != nil {
			return nil, err
		}

		var obj *domain.Object
		if err := json.Unmarshal(objectData, &obj); err != nil {
			return nil, err
		}

		if obj != nil {
			objects = append(objects, obj)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	result, err := s.db.ExecContext(ctx, `
		DELETE FROM objects
		WHERE id = ?
	`, id)
	
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
