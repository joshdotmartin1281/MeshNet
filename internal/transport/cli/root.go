package cli

import (
	"MeshNet/internal/app"
)

func Run(s *app.Service, args []string) error {
    if len(args) == 0 {
        return ErrNoCommand
    }

    switch args[0] {
    case "put":
        return Put(s, args[1:])
    //case "get":
        //return Get(s, args[1:])
    //case "list":
        //return List(s, args[1:])
    //case "delete":
        //return Delete(s, args[1:])
    default:
        return ErrUnknownCommand
    }
}

