package app

import (
	"context"
	"fmt"
	"net/url"

	"MeshNet/query"
)

type Queries struct {
	store   query.RelationalStore
	queries map[string]query.Query
}

func NewQueries(ctx context.Context, store query.RelationalStore, queries ...query.Query) (*Queries, error) {
	registry := make(map[string]query.Query)

	for _, qr := range queries {
		if err := qr.Init(ctx, store); err != nil {
			return nil, fmt.Errorf("init query %s: %w", qr.Name(), err)
		}

		if _, exists := registry[qr.Name()]; exists {
			return nil, fmt.Errorf("duplicate query name: %s", qr.Name())
		}

		registry[qr.Name()] = qr
	}

	return &Queries{
		store:   store,
		queries: registry,
	}, nil
}

func (q *Queries) Run(ctx context.Context, name string, params url.Values) (any, error) {
	qr, ok := q.queries[name]
	if !ok {
		return nil, fmt.Errorf("query not found: %s", name)
	}

	return qr.Run(ctx, q.store, params)
}

func (q *Queries) Mutate(ctx context.Context, name string, params url.Values) (any, error) {
	qr, ok := q.queries[name]
	if !ok {
		return nil, fmt.Errorf("query not found: %s", name)
	}

	mutable, ok := qr.(query.MutableQuery)
	if !ok {
		return nil, fmt.Errorf("query %s does not support write operations", name)
	}

	return mutable.Mutate(ctx, q.store, params)
}
