package app

import (
	"context"
	"net/url"

	"MeshNet/internal/api"
	"MeshNet/uauth"
)

type Port interface {
	ObjectPort
	RelationalPort
	AuthPort
	AdminPort
}

type AuthPort interface {
	AuthEnabled() bool
	Authenticate(ctx context.Context, token string) (auth.Identity, error)
}

type AdminPort interface {
	CreateUser(ctx context.Context, name string) (auth.User, error)
	ListUsers(ctx context.Context) ([]auth.User, error)
	IssueToken(ctx context.Context, user string) (string, error)
	RevokeTokens(ctx context.Context, user string) error
	Grant(ctx context.Context, user, collection string, level auth.Level) error
	Ungrant(ctx context.Context, user, collection string) error
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
