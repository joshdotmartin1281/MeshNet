package storage

import (
	"context"

	"MeshNet/internal/domain"
)

type ObjectStore interface {
	Save(ctx context.Context, obj *domain.Object, payload *domain.Payload) error

	Get(ctx context.Context, collection string, id string) (*domain.Object, *domain.Payload, error)

	GetByHash(ctx context.Context, collection string, hash string) (*domain.Object, *domain.Payload, error)

	List(ctx context.Context, collection string) ([]*domain.Object, error)

	Delete(ctx context.Context, collection string, id string) error

	DeleteByHash(ctx context.Context, collection string, hash string) error
}
