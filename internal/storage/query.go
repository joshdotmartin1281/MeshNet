package storage

import (
	"context"
	"net/url"
)

type Query interface {
	Name() string
	Init(ctx context.Context, store RelationalStore) error
	Run(ctx context.Context, store RelationalStore, params url.Values) (any, error)
}

type MutableQuery interface {
	Query
	Mutate(ctx context.Context, store RelationalStore, params url.Values) (any, error)
}
