package cli

import (
	"context"
	"errors"
	"flag"
	"io"

	"MeshNet/internal/app"
)

func Delete(s *app.Service, args []string) error {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New("usage: delete <id>")
	}

	id := fs.Arg(0)

	return s.Delete(context.Background(), id)
}

