package app

import (
	"context"
	"errors"
	"time"

	"MeshNet/internal/api"
	"MeshNet/internal/domain"
	"MeshNet/internal/storage"
)

type Service struct {
	repo      storage.ObjectStore
	hasher    domain.Hasher
	processor *Processor
}

func New(repo storage.ObjectStore, hasher domain.Hasher, processor *Processor) *Service {
	return &Service{
		repo:      repo,
		hasher:    hasher,
		processor: processor,
	}
}

func (s *Service) Put(ctx context.Context, req api.PutRequest) (api.PutResponse, error) {
	if err := ctx.Err(); err != nil {
		return api.PutResponse{}, err
	}

	if len(req.Data) == 0 {
		return api.PutResponse{}, errors.New("empty payload")
	}

	data := req.Data

	if len(req.Transforms) > 0 {
		var err error

		data, err = s.processor.Process(data, req.Transforms)
		if err != nil {
			return api.PutResponse{}, err
		}
	}

	if len(data) == 0 {
		return api.PutResponse{}, errors.New("empty payload after transforms")
	}

	if err := ctx.Err(); err != nil {
		return api.PutResponse{}, err
	}

	obj := &domain.Object{
		ID:         domain.NewID().String(),
		Collection: req.Collection,
		Name:       req.Name,
		MediaType:  req.MediaType,
		Size:       int64(len(data)),
		Source:     req.Source,
		CreatedAt:  time.Now().UTC(),
		Transforms: req.Transforms,
	}

	obj.Hash = s.hasher.Hash(data)

	payload := &domain.Payload{
		ObjectID: obj.ID,
		Data:     data,
	}

	if err := s.repo.Save(ctx, obj, payload); err != nil {
		return api.PutResponse{}, err
	}

	return api.PutResponse{
		Object: obj,
	}, nil
}

func (s *Service) Get(ctx context.Context, req api.GetRequest) (api.GetResponse, error) {
	obj, payload, err := s.repo.Get(ctx, req.Collection, req.ID)
	if err != nil {
		return api.GetResponse{}, err
	}

	data, err := s.processor.Process(payload.Data, req.Transforms)
	if err != nil {
		return api.GetResponse{}, err
	}

	payload.Data = data

	return api.GetResponse{
		Object:  obj,
		Payload: payload,
	}, nil
}

func (s *Service) GetByHash(ctx context.Context, req api.GetByHashRequest) (api.GetResponse, error) {
	obj, payload, err := s.repo.GetByHash(ctx, req.Collection, req.Hash)
	if err != nil {
		return api.GetResponse{}, err
	}

	data, err := s.processor.Process(payload.Data, req.Transforms)
	if err != nil {
		return api.GetResponse{}, err
	}

	payload.Data = data

	return api.GetResponse{
		Object:  obj,
		Payload: payload,
	}, nil
}

func (s *Service) List(ctx context.Context, req api.ListRequest) (api.ListResponse, error) {
	objects, err := s.repo.List(ctx, req.Collection)
	if err != nil {
		return api.ListResponse{}, err
	}

	return api.ListResponse{
		Objects: objects,
	}, nil
}

func (s *Service) Delete(ctx context.Context, req api.DeleteRequest) (api.DeleteResponse, error) {
	if err := s.repo.Delete(ctx, req.Collection, req.ID); err != nil {
		return api.DeleteResponse{}, err
	}

	return api.DeleteResponse{}, nil
}

func (s *Service) DeleteByHash(ctx context.Context, req api.DeleteByHashRequest) (api.DeleteByHashResponse, error) {
	if err := s.repo.DeleteByHash(ctx, req.Collection, req.Hash); err != nil {
		return api.DeleteByHashResponse{}, err
	}

	return api.DeleteByHashResponse{}, nil
}
