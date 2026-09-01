package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"MeshNet/internal/domain"
)

type record struct {
	Object  *domain.Object  `json:"object"`
	Payload *domain.Payload `json:"payload"`
}

type Store struct {
	root string
}

func New(root string) (*Store, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}

	return &Store{
		root: root,
	}, nil
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, "*.tmp")
	if err != nil {
		return err
	}

	tmpName := tmp.Name()

	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.Chmod(tmpName, perm); err != nil {
		return err
	}

	return os.Rename(tmpName, path)
}

func (s *Store) Save(ctx context.Context, obj *domain.Object, payload *domain.Payload) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := os.Stat(s.objectPath(obj.Collection, obj.ID)); err == nil {
		return domain.ErrDuplicate
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if _, err := os.Stat(s.hashPath(obj.Collection, obj.Hash)); err == nil {
		return domain.ErrDuplicate
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	rec := record{
		Object:  obj,
		Payload: payload,
	}

	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}

	if err := writeFileAtomic(s.objectPath(obj.Collection, obj.ID), data, 0644); err != nil {
		return err
	}

	if err := writeFileAtomic(s.hashPath(obj.Collection, obj.Hash), []byte(obj.ID), 0644); err != nil {
		_ = os.Remove(s.objectPath(obj.Collection, obj.ID))
		return err
	}

	return nil
}

func (s *Store) Get(ctx context.Context, collection string, id string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	data, err := os.ReadFile(s.objectPath(collection, id))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil, domain.ErrNotFound
		}

		return nil, nil, err
	}

	var rec record

	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, nil, err
	}

	return rec.Object, rec.Payload, nil
}

func (s *Store) GetByHash(ctx context.Context, collection string, hash string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	id, err := os.ReadFile(s.hashPath(collection, hash))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil, domain.ErrNotFound
		}

		return nil, nil, err
	}

	return s.Get(ctx, collection, string(id))
}

func (s *Store) List(ctx context.Context, collection string) ([]*domain.Object, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	dir := s.objectsDir(collection)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*domain.Object{}, nil
		}

		return nil, err
	}

	objects := make([]*domain.Object, 0, len(entries))

	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}

		var rec record

		if err := json.Unmarshal(data, &rec); err != nil {
			return nil, err
		}

		if rec.Object != nil {
			objects = append(objects, rec.Object)
		}
	}

	return objects, nil
}

func (s *Store) Delete(ctx context.Context, collection string, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	obj, _, err := s.Get(ctx, collection, id)
	if err != nil {
		return err
	}

	if err := os.Remove(s.objectPath(collection, obj.ID)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.Remove(s.hashPath(collection, obj.Hash)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func (s *Store) DeleteByHash(ctx context.Context, collection string, hash string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	obj, _, err := s.GetByHash(ctx, collection, hash)
	if err != nil {
		return err
	}

	return s.Delete(ctx, collection, obj.ID)
}

func (s *Store) objectsDir(collection string) string {
	return filepath.Join(s.root, "paths", collection, "objects")
}

func (s *Store) hashesDir(collection string) string {
	return filepath.Join(s.root, "paths", collection, "hashes")
}

func (s *Store) objectPath(collection string, id string) string {
	return filepath.Join(s.objectsDir(collection), id+".json")
}

func (s *Store) hashPath(collection string, hash string) string {
	return filepath.Join(s.hashesDir(collection), hash)
}
