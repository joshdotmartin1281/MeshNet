package app

import (
	"context"
	"time"

	"MeshNet/internal/domain"
	"MeshNet/internal/processors/hash"
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
	obj.Hash = hash.SHA256(payload.Data)
	obj.Size = int64(len(payload.Data))
	obj.CreatedAt = time.Now()

	payload.ObjectID = obj.ID

	if err := s.repo.SaveObject(ctx, obj); err != nil {
		return err
	}

	return s.repo.SavePayload(ctx, payload)
}

