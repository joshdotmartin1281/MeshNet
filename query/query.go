package query

import (
	"context"
	"database/sql"
	"net/url"
)

type RelationalStore interface {
	ExecContext(
		ctx context.Context,
		query string,
		args ...any,
	) (sql.Result, error)

	QueryContext(
		ctx context.Context,
		query string,
		args ...any,
	) (*sql.Rows, error)

	QueryRowContext(
		ctx context.Context,
		query string,
		args ...any,
	) *sql.Row
}

type Query interface {
	Name() string
	Init(ctx context.Context, store RelationalStore) error
	Run(ctx context.Context, store RelationalStore, params url.Values) (any, error)
}

type MutableQuery interface {
	Query
	Mutate(ctx context.Context, store RelationalStore, params url.Values) (any, error)
}
