package app

import (
	"context"
	"errors"
	"net/url"
	"time"

	"MeshNet/internal/api"
	"MeshNet/internal/domain"
	"MeshNet/internal/storage"
	"MeshNet/uauth"
)

type Service struct {
	repo      storage.ObjectStore
	hasher    domain.Hasher
	processor *Processor
	queries   *Queries
	auth      auth.Provider
	manager   auth.Manager
}

type Option func(*Service)

func WithAuth(p auth.Provider) Option {
	return func(s *Service) {
		s.auth = p

		if m, ok := p.(auth.Manager); ok {
			s.manager = m
		}
	}
}

func New(repo storage.ObjectStore, hasher domain.Hasher, processor *Processor, queries *Queries, opts ...Option) *Service {
	s := &Service{
		repo:      repo,
		hasher:    hasher,
		processor: processor,
		queries:   queries,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) AuthEnabled() bool {
	return s.auth != nil
}

func (s *Service) Authenticate(ctx context.Context, token string) (auth.Identity, error) {
	if s.auth == nil {
		return auth.Identity{}, auth.ErrNotSupported
	}

	id, err := s.auth.Authenticate(ctx, token)
	if err != nil {
		return auth.Identity{}, err
	}

	id.Admin = false

	return id, nil
}

func (s *Service) authorize(ctx context.Context, collection string, level auth.Level) error {
	if s.auth == nil {
		return nil
	}

	if id, ok := auth.FromContext(ctx); ok && id.Admin {
		return nil
	}

	return s.auth.Authorize(ctx, collection, level)
}

func (s *Service) admin(ctx context.Context) (auth.Manager, error) {
	if s.manager == nil {
		return nil, auth.ErrNotSupported
	}

	if id, ok := auth.FromContext(ctx); !ok || !id.Admin {
		return nil, auth.ErrForbidden
	}

	return s.manager, nil
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

func (s *Service) RunQuery(ctx context.Context, name string, params url.Values) (any, error) {
	return s.queries.Run(ctx, name, params)
}

func (s *Service) MutateQuery(ctx context.Context, name string, params url.Values) (any, error) {
	return s.queries.Mutate(ctx, name, params)
}

func (s *Service) CreateUser(ctx context.Context, name string) (auth.User, error) {
	m, err := s.admin(ctx)
	if err != nil {
		return auth.User{}, err
	}

	return m.CreateUser(ctx, name)
}

func (s *Service) ListUsers(ctx context.Context) ([]auth.User, error) {
	m, err := s.admin(ctx)
	if err != nil {
		return nil, err
	}

	return m.ListUsers(ctx)
}

func (s *Service) IssueToken(ctx context.Context, user string) (string, error) {
	m, err := s.admin(ctx)
	if err != nil {
		return "", err
	}

	return m.IssueToken(ctx, user)
}

func (s *Service) RevokeTokens(ctx context.Context, user string) error {
	m, err := s.admin(ctx)
	if err != nil {
		return err
	}

	return m.RevokeTokens(ctx, user)
}

func (s *Service) Grant(ctx context.Context, user, collection string, level auth.Level) error {
	m, err := s.admin(ctx)
	if err != nil {
		return err
	}

	return m.Grant(ctx, user, collection, level)
}

func (s *Service) Ungrant(ctx context.Context, user, collection string) error {
	m, err := s.admin(ctx)
	if err != nil {
		return err
	}

	return m.Ungrant(ctx, user, collection)
}
