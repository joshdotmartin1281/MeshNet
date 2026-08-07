package memory

import (
	"context"
	"errors"

	"MeshNet/internal/domain"
)

var (
	ErrNotFound  = errors.New("object not found")
	ErrDuplicate = errors.New("object already exists")
)

type Store struct {
	objects  map[string]*domain.Object
	payloads map[string]*domain.Payload
	hashes   map[string]string
}

func New() *Store {
	return &Store{
		objects:  make(map[string]*domain.Object),
		payloads: make(map[string]*domain.Payload),
		hashes:   make(map[string]string),
	}
}

func (s *Store) Save(ctx context.Context, obj *domain.Object, payload *domain.Payload) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, exists := s.hashes[obj.Hash]; exists {
		return ErrDuplicate
	}

	s.objects[obj.ID] = obj
	s.payloads[obj.ID] = payload
	s.hashes[obj.Hash] = obj.ID

	return nil
}

func (s *Store) Get(ctx context.Context, id string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	obj, ok := s.objects[id]
	if !ok {
		return nil, nil, ErrNotFound
	}

	payload, ok := s.payloads[id]
	if !ok {
		return nil, nil, ErrNotFound
	}

	return obj, payload, nil
}

func (s *Store) GetByHash(ctx context.Context, hash string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	id, ok := s.hashes[hash]
	if !ok {
		return nil, nil, ErrNotFound
	}

	return s.Get(ctx, id)
}

func (s *Store) List(ctx context.Context) ([]*domain.Object, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	objects := make([]*domain.Object, 0, len(s.objects))

	for _, obj := range s.objects {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		objects = append(objects, obj)
	}

	return objects, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	obj, ok := s.objects[id]
	if !ok {
		return ErrNotFound
	}

	delete(s.objects, id)
	delete(s.payloads, id)
	delete(s.hashes, obj.Hash)

	return nil
}

