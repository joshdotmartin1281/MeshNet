package domain

import (
	"context"
	"errors"
)

var (
	ErrNotFound  = errors.New("object not found")
	ErrDuplicate = errors.New("object already exists")
)

type ObjectRepository interface {
	Save(ctx context.Context, obj *Object, payload *Payload) error

	Get(ctx context.Context, collection string, id string) (*Object, *Payload, error)

	GetByHash(ctx context.Context, collection string, hash string) (*Object, *Payload, error)

	List(ctx context.Context, collection string) ([]*Object, error)

	Delete(ctx context.Context, collection string, id string) error

	DeleteByHash(ctx context.Context, collection string, hash string) error
}
