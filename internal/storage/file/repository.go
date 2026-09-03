package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"MeshNet/internal/domain"
)

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

func writeFileAtomic(collection string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(collection)

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

	return os.Rename(tmpName, collection)
}

func (s *Store) Save(ctx context.Context, obj *domain.Object, payload *domain.Payload) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	objectPath := s.objectPath(obj.Collection, obj.ID)
	hashPath := s.hashPath(obj.Collection, obj.Hash)
	payloadPath := s.payloadPath(obj.Collection, obj.Hash)

	if _, err := os.Stat(objectPath); err == nil {
		return domain.ErrDuplicate
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if _, err := os.Stat(hashPath); err == nil {
		return domain.ErrDuplicate
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if payload != nil {
		if payload.ObjectID != "" && payload.ObjectID != obj.ID {
			return errors.New("payload ObjectID does not match object ID")
		}

		if int64(len(payload.Data)) != obj.Size {
			return errors.New("payload size does not match object size")
		}
	}

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	if err := writeFileAtomic(objectPath, data, 0644); err != nil {
		return err
	}

	if payload != nil {
		if err := writeFileAtomic(payloadPath, payload.Data, 0644); err != nil {
			_ = os.Remove(objectPath)
			return err
		}
	}

	if err := writeFileAtomic(hashPath, []byte(obj.ID), 0644); err != nil {
		_ = os.Remove(objectPath)
		_ = os.Remove(payloadPath)
		return err
	}

	return nil
}

func (s *Store) getObject(ctx context.Context, collection string, id string) (*domain.Object, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(s.objectPath(collection, id))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	var obj domain.Object

	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

func (s *Store) Get(ctx context.Context, collection string, id string) (*domain.Object, *domain.Payload, error) {
	obj, err := s.getObject(ctx, collection, id)
	if err != nil {
		return nil, nil, err
	}

	payloadData, err := os.ReadFile(s.payloadPath(collection, obj.Hash))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return obj, nil, domain.ErrNotFound
		}

		return nil, nil, err
	}

	return obj, &domain.Payload{
		ObjectID: obj.ID,
		Data:     payloadData,
	}, nil
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

		var obj domain.Object

		if err := json.Unmarshal(data, &obj); err != nil {
			return nil, err
		}

		objects = append(objects, &obj)
	}

	return objects, nil
}

func (s *Store) removeEmptyCollectionDirs(collection string) {
	_ = os.Remove(s.objectsDir(collection))
	_ = os.Remove(s.hashesDir(collection))
	_ = os.Remove(s.collectionDir(collection))
}

func (s *Store) Delete(ctx context.Context, collection string, id string) error {
	obj, err := s.getObject(ctx, collection, id)
	if err != nil {
		return err
	}

	if err := os.Remove(s.objectPath(collection, obj.ID)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.Remove(s.hashPath(collection, obj.Hash)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.Remove(s.payloadPath(collection, obj.Hash)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	s.removeEmptyCollectionDirs(collection)

	return nil
}

func (s *Store) DeleteByHash(ctx context.Context, collection string, hash string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	id, err := os.ReadFile(s.hashPath(collection, hash))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.ErrNotFound
		}

		return err
	}

	return s.Delete(ctx, collection, string(id))
}

func (s *Store) objectsDir(collection string) string {
	return filepath.Join(s.root, collection, "objects")
}

func (s *Store) hashesDir(collection string) string {
	return filepath.Join(s.root, collection, "hashes")
}

func (s *Store) collectionDir(collection string) string {
	return filepath.Join(s.root, collection) 
}

func (s *Store) objectPath(collection string, id string) string {
	return filepath.Join(s.objectsDir(collection), id+".json")
}

func (s *Store) hashPath(collection string, hash string) string {
	return filepath.Join(s.hashesDir(collection), hash)
}

func (s *Store) payloadPath(collection string, hash string) string {
	return filepath.Join(s.root, collection, hash+".bin")
}
