package cli

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"

	"MeshNet/internal/app"
)

func Get(s *app.Service, args []string) error {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)

	out := fs.String("o", "", "output file")
	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New("usage: get [-o file] <id>")
	}

	id := fs.Arg(0)

	_, payload, err := s.Get(context.Background(), id)
	if err != nil {
		return err
	}

	if *out != "" {
		return os.WriteFile(*out, payload.Data, 0o644)
	}

	_, err = os.Stdout.Write(payload.Data)
	return err
}

func GetByHash(s *app.Service, args []string) error {
	fs := flag.NewFlagSet("get-by-hash", flag.ContinueOnError)
	out := fs.String("o", "", "output file")
	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New("usage: get-by-hash [-o file] <hash>")
	}

	hash := fs.Arg(0)

	_, payload, err := s.GetByHash(context.Background(), hash)
	if err != nil {
		return err
	}

	if *out != "" {
		return os.WriteFile(*out, payload.Data, 0o644)
	}

	_, err = os.Stdout.Write(payload.Data)
	return err
}

