package file

import (
	"os"
	"errors"
	"context"
	"path/filepath"
	"encoding/json"

	"MeshNet/internal/domain"
)

var ErrNotFound = errors.New("object not found")

type record struct {
    Object  *domain.Object  `json:"object"`
    Payload *domain.Payload `json:"payload"`
}

type Store struct {
    root string
}

func New(root string) (*Store, error) {
	objectsDir := filepath.Join(root, "objects")

	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		return nil, err
	}

	return &Store{
		root: root,
	}, nil
}

func (s *Store) Save(ctx context.Context, obj *domain.Object, payload *domain.Payload) error {
	record := record {
		Object: obj,
		Payload: payload,
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		s.objectPath(obj.ID),
		data,
		0644,
	)
}

func (s *Store) Get(ctx context.Context, id string) (*domain.Object, *domain.Payload, error) {
	path := s.objectPath(id)

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil, ErrNotFound
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
	return nil, nil, nil
}

func (s *Store) List(ctx context.Context) ([]*domain.Object, error) {
	return nil, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *Store) objectPath(id string) string {
    return filepath.Join(s.root, "objects", id+".json")
}

