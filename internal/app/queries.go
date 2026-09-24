package app

import (
	"context"
	"fmt"
	"net/url"

	"MeshNet/internal/storage"
)

type Queries struct {
	store   storage.RelationalStore
	queries map[string]storage.Query
}

func NewQueries(ctx context.Context, store storage.RelationalStore, queries ...storage.Query) (*Queries, error) {
	registry := make(map[string]storage.Query)

	for _, query := range queries {
		if err := query.Init(ctx, store); err != nil {
			return nil, fmt.Errorf("init query %s: %w", query.Name(), err)
		}

		if _, exists := registry[query.Name()]; exists {
			return nil, fmt.Errorf("duplicate query name: %s", query.Name())
		}

		registry[query.Name()] = query
	}

	return &Queries{
		store:   store,
		queries: registry,
	}, nil
}

func (q *Queries) Run(ctx context.Context, name string, params url.Values) (any, error) {
	query, ok := q.queries[name]
	if !ok {
		return nil, fmt.Errorf("query not found: %s", name)
	}

	return query.Run(ctx, q.store, params)
}

func (q *Queries) Mutate(ctx context.Context, name string, params url.Values) (any, error) {
	query, ok := q.queries[name]
	if !ok {
		return nil, fmt.Errorf("query not found: %s", name)
	}

	mutable, ok := query.(storage.MutableQuery)
	if !ok {
		return nil, fmt.Errorf("query %s does not support write operations", name)
	}

	return mutable.Mutate(ctx, q.store, params)
}
