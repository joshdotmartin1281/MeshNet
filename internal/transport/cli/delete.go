package cli

import (
	"context"
	"errors"
	"flag"
	"io"

	"MeshNet/internal/api"
	"MeshNet/internal/app"
)

func Delete(port app.Port, args []string) error {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New("usage: delete <id>")
	}

	_, err := port.Delete(
		context.Background(),
		api.DeleteRequest{
			ID: fs.Arg(0),
		},
	)

	return err
}

