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

		CREATE INDEX IF NOT EXISTS idx_objects_hash 
			ON objects(hash);
		`)
	return err
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

func isUniqueViolation(err error) bool {
	var sqliteErr interface {
		ErrorCode() int
	}

	if !errors.As(err, &sqliteErr) {
		return false
	}

	// SQLITE_CONSTRAINT = 19.
	return sqliteErr.ErrorCode() == 19
}

func (s *Store) Get(
	ctx context.Context,
	id string,
) (*domain.Object, *domain.Payload, error) {
	panic("not implemented")
}

func (s *Store) GetByHash(
	ctx context.Context,
	hash string,
) (*domain.Object, *domain.Payload, error) {
	panic("not implemented")
}

func (s *Store) List(
	ctx context.Context,
) ([]*domain.Object, error) {
	panic("not implemented")
}

func (s *Store) Delete(
	ctx context.Context,
	id string,
) error {
	panic("not implemented")
}
