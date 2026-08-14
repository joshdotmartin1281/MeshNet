package app

import (
	"context"

	"MeshNet/internal/api"
)

type Port interface {
	Put(context.Context, api.PutRequest) (api.PutResponse, error)
	Get(context.Context, api.GetRequest) (api.GetResponse, error)
	GetByHash(context.Context, api.GetByHashRequest) (api.GetResponse, error)
	List(context.Context, api.ListRequest) (api.ListResponse, error)
	Delete(context.Context, api.DeleteRequest) (api.DeleteResponse, error)
}

