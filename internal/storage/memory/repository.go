package memory

import (
	"context"

	"MeshNet/internal/domain"
)

type Store struct {
	objects map[string]*domain.Object
	payloads map[string]*domain.Payload
}

func New() *Store {
	return &Store{
		objects: make(map[string]*domain.Object),
		payloads: make(map[string]*domain.Payload),
	}
}

func (s *Store) SaveObject(ctx context.Context, obj *domain.Object) error {
	s.objects[obj.ID] = obj
	return nil
}

func (s *Store) SavePayload(ctx context.Context, payload *domain.Payload) error {
	s.payloads[payload.ObjectID] = payload
	return nil
}

