package cli

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"

	"MeshNet/internal/api"
	"MeshNet/internal/app"
)

func Get(port app.Port, args []string) error {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)

	out := fs.String("o", "", "output file")
	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New("usage: get [-o file] <id>")
	}

	resp, err := port.Get(
		context.Background(),
		api.GetRequest{
			ID: fs.Arg(0),
		},
	)
	if err != nil {
		return err
	}

	if *out != "" {
		return os.WriteFile(*out, resp.Payload.Data, 0o644)
	}

	_, err = os.Stdout.Write(resp.Payload.Data)
	return err
}

func GetByHash(port app.Port, args []string) error {
	fs := flag.NewFlagSet("get-by-hash", flag.ContinueOnError)

	out := fs.String("o", "", "output file")
	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New("usage: get-by-hash [-o file] <hash>")
	}

	resp, err := port.GetByHash(
		context.Background(),
		api.GetByHashRequest{
			Hash: fs.Arg(0),
		},
	)
	if err != nil {
		return err
	}

	if *out != "" {
		return os.WriteFile(*out, resp.Payload.Data, 0o644)
	}

	_, err = os.Stdout.Write(resp.Payload.Data)
	return err
}

