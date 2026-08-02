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

func (s *Service) Process(ctx context.Context, source domain.Source, data []byte) (*domain.Object, error) {
	obj := &domain.Object{
		ID:        domain.NewID().String(),
		Source:    source,
		Hash:      hash.SHA256(data),
		Size:      int64(len(data)),
		CreatedAt: time.Now().UTC(),
	}

	payload := &domain.Payload{
		ObjectID: obj.ID,
		Data:     data,
	}

	if err := s.repo.Save(ctx, obj, payload); err != nil {
		return nil, err
	}

	return obj, nil
}

