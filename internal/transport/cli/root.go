package cli

import "MeshNet/internal/app"

func Run(port app.Port, args []string) error {
	if len(args) == 0 {
		return ErrNoCommand
	}

	switch args[0] {
	case "put":
		return Put(port, args[1:])
	case "get":
		return Get(port, args[1:])
	case "list":
		return List(port, args[1:])
	case "delete":
		return Delete(port, args[1:])
	case "get-by-hash":
		return GetByHash(port, args[1:])
	default:
		return ErrUnknownCommand
	}
}

