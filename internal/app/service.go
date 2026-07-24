package app

import (
	"context"

	"MeshNet/internal/domain"
)

type Service struct {
	repo domain.Repository
}

func New(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Process(ctx context.Context, msg *domain.Message) error {
	return s.repo.Save(ctx, msg)
}

