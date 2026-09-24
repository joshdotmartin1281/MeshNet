package app

import (
	"context"
	"net/url"

	"MeshNet/internal/api"
)

type Port interface {
	ObjectPort
	RelationalPort
}

type ObjectPort interface {
	Put(context.Context, api.PutRequest) (api.PutResponse, error)
	Get(context.Context, api.GetRequest) (api.GetResponse, error)
	GetByHash(context.Context, api.GetByHashRequest) (api.GetResponse, error)
	List(context.Context, api.ListRequest) (api.ListResponse, error)
	Delete(context.Context, api.DeleteRequest) (api.DeleteResponse, error)
	DeleteByHash(context.Context, api.DeleteByHashRequest) (api.DeleteByHashResponse, error)
}

type RelationalPort interface {
	RunQuery(ctx context.Context, name string, params url.Values) (any, error)
	MutateQuery(ctx context.Context, name string, params url.Values) (any, error)
}
