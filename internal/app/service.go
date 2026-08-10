package app

import (
	"context"
	"errors"
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

func (s *Service) Put(ctx context.Context, source domain.Source, data []byte) (*domain.Object, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("empty payload")
	}

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

func (s *Service) Get(ctx context.Context, id string) (*domain.Object, *domain.Payload, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetByHash(ctx context.Context, hash string) (*domain.Object, *domain.Payload, error) {
	return s.repo.GetByHash(ctx, hash)
}

func (s *Service) List(ctx context.Context) ([]*domain.Object, error) {
	return s.repo.List(ctx)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

