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

	Get(ctx context.Context, id string) (*Object, *Payload, error)

	GetByHash(ctx context.Context, hash string) (*Object, *Payload, error)

	List(ctx context.Context) ([]*Object, error)

	Delete(ctx context.Context, id string) error
}
