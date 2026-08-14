package app

import (
	"context"
	"errors"
	"time"

	"MeshNet/internal/api"
	"MeshNet/internal/domain"
)

type Service struct {
	repo   domain.ObjectRepository
	hasher domain.Hasher
}

func New(repo domain.ObjectRepository, hasher domain.Hasher) *Service {
	return &Service{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *Service) Put(ctx context.Context, req api.PutRequest) (api.PutResponse, error) {
	if err := ctx.Err(); err != nil {
		return api.PutResponse{}, err
	}

	if len(req.Data) == 0 {
		return api.PutResponse{}, errors.New("empty payload")
	}

	obj := &domain.Object{
		ID:        domain.NewID().String(),
		Source:    req.Source,
		Size:      int64(len(req.Data)),
		CreatedAt: time.Now().UTC(),
	}

	obj.Hash = s.hasher.Hash(req.Data)

	payload := &domain.Payload{
		ObjectID: obj.ID,
		Data:     req.Data,
	}

	if err := s.repo.Save(ctx, obj, payload); err != nil {
		return api.PutResponse{}, err
	}

	return api.PutResponse{
		Object: obj,
	}, nil
}

func (s *Service) Get(ctx context.Context, req api.GetRequest) (api.GetResponse, error) {
	obj, payload, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		return api.GetResponse{}, err
	}

	return api.GetResponse{
		Object:  obj,
		Payload: payload,
	}, nil
}

func (s *Service) GetByHash(ctx context.Context, req api.GetByHashRequest) (api.GetResponse, error) {
	obj, payload, err := s.repo.GetByHash(ctx, req.Hash)
	if err != nil {
		return api.GetResponse{}, err
	}

	return api.GetResponse{
		Object:  obj,
		Payload: payload,
	}, nil
}

func (s *Service) List(ctx context.Context, req api.ListRequest) (api.ListResponse, error) {
	objects, err := s.repo.List(ctx)
	if err != nil {
		return api.ListResponse{}, err
	}

	return api.ListResponse{
		Objects: objects,
	}, nil
}

func (s *Service) Delete(ctx context.Context, req api.DeleteRequest) (api.DeleteResponse, error) {
	if err := s.repo.Delete(ctx, req.ID); err != nil {
		return api.DeleteResponse{}, err
	}

	return api.DeleteResponse{}, nil
}

