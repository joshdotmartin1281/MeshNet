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
	objectsDir := filepath.Join(root, "objects")
	hashesDir := filepath.Join(root, "hashes")

	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(hashesDir, 0755); err != nil {
		return nil, err
	}

	return &Store{
		root: root,
	}, nil
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)

	tmp, err := os.CreateTemp(dir, "*.tmp")
	if err != nil {
		return err
	}

	tmpName := tmp.Name()

	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}

	if err := tmp.Sync(); err != nil {
		tmp.Close()
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

	if _, err := os.Stat(s.objectPath(obj.ID)); err == nil {
		return domain.ErrDuplicate
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if _, err := os.Stat(s.hashPath(obj.Hash)); err == nil {
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

	if err := writeFileAtomic(
		s.objectPath(obj.ID),
		data,
		0644,
	); err != nil {
		return err
	}

	if err := writeFileAtomic(
		s.hashPath(obj.Hash),
		[]byte(obj.ID),
		0644,
	); err != nil {
		_ = os.Remove(s.objectPath(obj.ID))
		return err
	}

	return nil
}

func (s *Store) Get(ctx context.Context, id string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	data, err := os.ReadFile(s.objectPath(id))
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

func (s *Store) GetByHash(ctx context.Context, hash string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	id, err := os.ReadFile(s.hashPath(hash))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil, domain.ErrNotFound
		}

		return nil, nil, err
	}

	return s.Get(ctx, string(id))
}

func (s *Store) List(ctx context.Context) ([]*domain.Object, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	dir := filepath.Join(s.root, "objects")

	entries, err := os.ReadDir(dir)
	if err != nil {
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

		data, err := os.ReadFile(
			filepath.Join(dir, entry.Name()),
		)
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

func (s *Store) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	obj, _, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := os.Remove(s.objectPath(id)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.ErrNotFound
		}

		return err
	}

	if err := os.Remove(s.hashPath(obj.Hash)); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func (s *Store) objectPath(id string) string {
	return filepath.Join(s.root, "objects", id+".json")
}

func (s *Store) hashPath(hash string) string {
	return filepath.Join(s.root, "hashes", hash)
}
