package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"MeshNet/internal/domain"
	"MeshNet/internal/storage"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

var _ storage.ObjectStore = (*Store)(nil)
var _ storage.RelationalStore = (*Store)(nil)

func New(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

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

func (s *Store) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *Store) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *Store) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *Store) init() error {
	_, err := s.db.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS collections (
			id         TEXT PRIMARY KEY,
			created_at TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS objects (
			id            TEXT PRIMARY KEY,
			collection_id TEXT NOT NULL,
			hash          TEXT NOT NULL,
			name          TEXT NOT NULL,
			media_type    TEXT NOT NULL,
			size          INTEGER NOT NULL,
			source        BLOB NOT NULL,
			created_at    TEXT NOT NULL,
			transforms    BLOB NOT NULL,

			FOREIGN KEY (collection_id)
				REFERENCES collections(id)
				ON DELETE CASCADE,

			UNIQUE (collection_id, hash)
		);

		CREATE TABLE IF NOT EXISTS payloads (
			object_id TEXT PRIMARY KEY,
			data      BLOB NOT NULL,

			FOREIGN KEY (object_id)
				REFERENCES objects(id)
				ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_objects_collection
			ON objects(collection_id);

		CREATE INDEX IF NOT EXISTS idx_objects_hash
			ON objects(hash);
	`)

	return err
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func (s *Store) Save(ctx context.Context, obj *domain.Object, payload *domain.Payload) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if obj == nil {
		return errors.New("object is nil")
	}

	if payload != nil {
		if payload.ObjectID != "" && payload.ObjectID != obj.ID {
			return errors.New("payload ObjectID does not match object ID")
		}

		if int64(len(payload.Data)) != obj.Size {
			return errors.New("payload size does not match object size")
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	createdAt := obj.CreatedAt.Format(time.RFC3339Nano)

	_, err = tx.ExecContext(ctx, `
		INSERT INTO collections (
			id,
			created_at
		)
		VALUES (?, ?)
		ON CONFLICT(id) DO NOTHING
	`, obj.Collection, createdAt)
	if err != nil {
		return err
	}

	sourceData, err := json.Marshal(obj.Source)
	if err != nil {
		return err
	}

	transformsData, err := json.Marshal(obj.Transforms)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO objects (
			id,
			collection_id,
			hash,
			name,
			media_type,
			size,
			source,
			created_at,
			transforms
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		obj.ID,
		obj.Collection,
		obj.Hash,
		obj.Name,
		obj.MediaType,
		obj.Size,
		sourceData,
		createdAt,
		transformsData,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrDuplicate
		}

		return err
	}

	if payload != nil {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO payloads (
				object_id,
				data
			)
			VALUES (?, ?)
		`, obj.ID, payload.Data)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) Get(ctx context.Context, collection string, id string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	var (
		hash           string
		name           string
		mediaType      string
		size           int64
		sourceData     []byte
		createdAtData  string
		transformsData []byte
	)

	err := s.db.QueryRowContext(ctx, `
		SELECT
			hash,
			name,
			media_type,
			size,
			source,
			created_at,
			transforms
		FROM objects
		WHERE id = ?
		  AND collection_id = ?
	`, id, collection).Scan(
		&hash,
		&name,
		&mediaType,
		&size,
		&sourceData,
		&createdAtData,
		&transformsData,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, domain.ErrNotFound
		}

		return nil, nil, err
	}

	source, err := decodeSource(sourceData)
	if err != nil {
		return nil, nil, err
	}

	transforms, err := decodeTransforms(transformsData)
	if err != nil {
		return nil, nil, err
	}

	createdAt, err := time.Parse(time.RFC3339Nano, createdAtData)
	if err != nil {
		return nil, nil, err
	}

	obj := &domain.Object{
		ID:         id,
		Hash:       hash,
		Collection: collection,
		Name:       name,
		MediaType:  mediaType,
		Size:       size,
		Source:     source,
		CreatedAt:  createdAt,
		Transforms: transforms,
	}

	var payloadData []byte

	err = s.db.QueryRowContext(ctx, `
		SELECT data
		FROM payloads
		WHERE object_id = ?
	`, id).Scan(&payloadData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return obj, nil, nil
		}

		return nil, nil, err
	}

	payload := &domain.Payload{
		ObjectID: obj.ID,
		Data:     payloadData,
	}

	return obj, payload, nil
}

func (s *Store) GetByHash(ctx context.Context, collection string, hash string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	var id string

	err := s.db.QueryRowContext(ctx, `
		SELECT id
		FROM objects
		WHERE collection_id = ?
		  AND hash = ?
	`, collection, hash).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, domain.ErrNotFound
		}

		return nil, nil, err
	}

	return s.Get(ctx, collection, id)
}

func (s *Store) List(ctx context.Context, collection string) ([]*domain.Object, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT
			id,
			hash,
			name,
			media_type,
			size,
			source,
			created_at,
			transforms
		FROM objects
		WHERE collection_id = ?
		ORDER BY created_at ASC
	`, collection)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := make([]*domain.Object, 0)

	for rows.Next() {
		var (
			id             string
			hash           string
			name           string
			mediaType      string
			size           int64
			sourceData     []byte
			createdAtData  string
			transformsData []byte
		)

		if err := rows.Scan(
			&id,
			&hash,
			&name,
			&mediaType,
			&size,
			&sourceData,
			&createdAtData,
			&transformsData,
		); err != nil {
			return nil, err
		}

		source, err := decodeSource(sourceData)
		if err != nil {
			return nil, err
		}

		transforms, err := decodeTransforms(transformsData)
		if err != nil {
			return nil, err
		}

		createdAt, err := time.Parse(time.RFC3339Nano, createdAtData)
		if err != nil {
			return nil, err
		}

		objects = append(objects, &domain.Object{
			ID:         id,
			Hash:       hash,
			Collection: collection,
			Name:       name,
			MediaType:  mediaType,
			Size:       size,
			Source:     source,
			CreatedAt:  createdAt,
			Transforms: transforms,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func (s *Store) Delete(ctx context.Context, collection string, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	result, err := s.db.ExecContext(ctx, `
		DELETE FROM objects
		WHERE id = ?
		  AND collection_id = ?
	`, id, collection)
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

func (s *Store) DeleteByHash(ctx context.Context, collection string, hash string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	result, err := s.db.ExecContext(ctx, `
		DELETE FROM objects
		WHERE collection_id = ?
		  AND hash = ?
	`, collection, hash)
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

func decodeSource(data []byte) (domain.Source, error) {
	var source domain.Source

	if err := json.Unmarshal(data, &source); err != nil {
		return source, err
	}

	return source, nil
}

func decodeTransforms(data []byte) ([]domain.Transform, error) {
	var transforms []domain.Transform

	if err := json.Unmarshal(data, &transforms); err != nil {
		return nil, err
	}

	return transforms, nil
}
