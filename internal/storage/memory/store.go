package memory

import (
	"context"

	"MeshNet/internal/domain"
	"MeshNet/internal/storage"
)

type Store struct {
	objects  map[string]map[string]*domain.Object
	payloads map[string]map[string]*domain.Payload
	hashes   map[string]map[string]string
}

var _ storage.ObjectStore = (*Store)(nil)

func New() *Store {
	return &Store{
		objects:  make(map[string]map[string]*domain.Object),
		payloads: make(map[string]map[string]*domain.Payload),
		hashes:   make(map[string]map[string]string),
	}
}

func (s *Store) Save(
	ctx context.Context,
	obj *domain.Object,
	payload *domain.Payload,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	collection := obj.Collection

	if _, ok := s.objects[collection]; !ok {
		s.objects[collection] = make(map[string]*domain.Object)
		s.payloads[collection] = make(map[string]*domain.Payload)
		s.hashes[collection] = make(map[string]string)
	}

	if _, exists := s.objects[collection][obj.ID]; exists {
		return domain.ErrDuplicate
	}

	if _, exists := s.hashes[collection][obj.Hash]; exists {
		return domain.ErrDuplicate
	}

	s.objects[collection][obj.ID] = obj
	s.payloads[collection][obj.ID] = payload
	s.hashes[collection][obj.Hash] = obj.ID

	return nil
}

func (s *Store) Get(
	ctx context.Context,
	collection string,
	id string,
) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	objects, ok := s.objects[collection]
	if !ok {
		return nil, nil, domain.ErrNotFound
	}

	obj, ok := objects[id]
	if !ok {
		return nil, nil, domain.ErrNotFound
	}

	payload, ok := s.payloads[collection][id]
	if !ok {
		return nil, nil, domain.ErrNotFound
	}

	return obj, payload, nil
}

func (s *Store) GetByHash(
	ctx context.Context,
	collection string,
	hash string,
) (*domain.Object, *domain.Payload, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	hashes, ok := s.hashes[collection]
	if !ok {
		return nil, nil, domain.ErrNotFound
	}

	id, ok := hashes[hash]
	if !ok {
		return nil, nil, domain.ErrNotFound
	}

	return s.Get(ctx, collection, id)
}

func (s *Store) List(
	ctx context.Context,
	collection string,
) ([]*domain.Object, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	objectsByCollection, ok := s.objects[collection]
	if !ok {
		return []*domain.Object{}, nil
	}

	objects := make([]*domain.Object, 0, len(objectsByCollection))

	for _, obj := range objectsByCollection {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		objects = append(objects, obj)
	}

	return objects, nil
}

func (s *Store) Delete(
	ctx context.Context,
	collection string,
	id string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	objects, ok := s.objects[collection]
	if !ok {
		return domain.ErrNotFound
	}

	obj, ok := objects[id]
	if !ok {
		return domain.ErrNotFound
	}

	delete(s.objects[collection], id)
	delete(s.payloads[collection], id)
	delete(s.hashes[collection], obj.Hash)

	return nil
}

func (s *Store) DeleteByHash(
	ctx context.Context,
	collection string,
	hash string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	hashes, ok := s.hashes[collection]
	if !ok {
		return domain.ErrNotFound
	}

	id, ok := hashes[hash]
	if !ok {
		return domain.ErrNotFound
	}

	return s.Delete(ctx, collection, id)
}
