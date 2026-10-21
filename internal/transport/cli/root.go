package cli

import (
	"context"

	"MeshNet/internal/app"
)

func Run(ctx context.Context, port app.Port, args []string) error {
	if len(args) == 0 {
		return ErrNoCommand
	}

	switch args[0] {
	case "put":
		return Put(ctx, port, args[1:])
	case "get":
		return Get(ctx, port, args[1:])
	case "list":
		return List(ctx, port, args[1:])
	case "delete":
		return Delete(ctx, port, args[1:])
	case "delete-by-hash":
		return DeleteByHash(ctx, port, args[1:])
	case "get-by-hash":
		return GetByHash(ctx, port, args[1:])
	case "query":
		return Query(ctx, port, args[1:])
	case "mutate-query":
		return MutateQuery(ctx, port, args[1:])
	case "user-create":
		return UserCreate(ctx, port, args[1:])
	case "user-list":
		return UserList(ctx, port, args[1:])
	case "token-create":
		return TokenCreate(ctx, port, args[1:])
	case "token-revoke":
		return TokenRevoke(ctx, port, args[1:])
	case "grant":
		return Grant(ctx, port, args[1:])
	case "ungrant":
		return Ungrant(ctx, port, args[1:])
	default:
		return ErrUnknownCommand
	}
}
