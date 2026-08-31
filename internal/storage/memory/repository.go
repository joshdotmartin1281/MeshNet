package memory

import (
	"context"

	"MeshNet/internal/domain"
)

type Store struct {
	objects  map[string]map[string]*domain.Object
	payloads map[string]map[string]*domain.Payload
	hashes   map[string]map[string]string
}

func New() *Store {
	return &Store{
		objects:  make(map[string]map[string]*domain.Object),
		payloads: make(map[string]map[string]*domain.Payload),
		hashes:   make(map[string]map[string]string),
	}
}

func (s *Store) Save(ctx context.Context, obj *domain.Object, payload *domain.Payload) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	path := obj.Path

	if _, ok := s.objects[path]; !ok {
		s.objects[path] = make(map[string]*domain.Object)
		s.payloads[path] = make(map[string]*domain.Payload)
		s.hashes[path] = make(map[string]string)
	}

	if _, exists := s.objects[path][obj.ID]; exists {
		return domain.ErrDuplicate
	}

	if _, exists := s.hashes[path][obj.Hash]; exists {
		return domain.ErrDuplicate
	}

	s.objects[path][obj.ID] = obj
	s.payloads[path][obj.ID] = payload
	s.hashes[path][obj.Hash] = obj.ID

	return nil
}

func (s *Store) Get(ctx context.Context, path string, id string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	objects, ok := s.objects[path]
	if !ok {
		return nil, nil, domain.ErrNotFound
	}

	obj, ok := objects[id]
	if !ok {
		return nil, nil, domain.ErrNotFound
	}

	payload, ok := s.payloads[path][id]
	if !ok {
		return nil, nil, domain.ErrNotFound
	}

	return obj, payload, nil
}

func (s *Store) GetByHash(ctx context.Context, path string, hash string) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	hashes, ok := s.hashes[path]
	if !ok {
		return nil, nil, domain.ErrNotFound
	}

	id, ok := hashes[hash]
	if !ok {
		return nil, nil, domain.ErrNotFound
	}

	return s.Get(ctx, path, id)
}

func (s *Store) List(ctx context.Context, path string) ([]*domain.Object, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	objectsByPath, ok := s.objects[path]
	if !ok {
		return []*domain.Object{}, nil
	}

	objects := make([]*domain.Object, 0, len(objectsByPath))

	for _, obj := range objectsByPath {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		objects = append(objects, obj)
	}

	return objects, nil
}

func (s *Store) Delete(ctx context.Context, path string, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	objects, ok := s.objects[path]
	if !ok {
		return domain.ErrNotFound
	}

	obj, ok := objects[id]
	if !ok {
		return domain.ErrNotFound
	}

	delete(s.objects[path], id)
	delete(s.payloads[path], id)
	delete(s.hashes[path], obj.Hash)

	return nil
}

func (s *Store) DeleteByHash(ctx context.Context, path string, hash string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	hashes, ok := s.hashes[path]
	if !ok {
		return domain.ErrNotFound
	}

	id, ok := hashes[hash]
	if !ok {
		return domain.ErrNotFound
	}

	return s.Delete(ctx, path, id)
}
