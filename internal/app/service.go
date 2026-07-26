package app

import (
	"context"
	"time"

	"MeshNet/internal/domain"
)

type Service struct {
	repo domain.ObjectRepository
}

func New(repo domain.ObjectRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Process(ctx context.Context, payload *domain.Payload, obj *domain.Object) error {
	// Calculate hash

	obj.Size = int64(len(payload.Data))
	obj.CreatedAt = time.Now()

	payload.ObjectID = obj.ID

	if err := s.repo.SaveObject(ctx, obj); err != nil {
		return err
	}

	return s.repo.SavePayload(ctx, payload)
}

