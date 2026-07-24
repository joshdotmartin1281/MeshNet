package memory

import (
	"context"

	"MeshNet/internal/domain"
)

type Repository struct {
	data []*domain.Message
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) Save(ctx context.Context, msg *domain.Message) error {
	r.data = append(r.data, msg)
	return nil
}

