package cli

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"

	"MeshNet/internal/api"
	"MeshNet/internal/app"
	"MeshNet/internal/domain"
)

func Get(ctx context.Context, port app.ObjectPort, args []string) error {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)

	out := fs.String("o", "", "output file")
	collection := fs.String("c", "", "collection to retrieve object")

	var transformArgs []string

	fs.Func(
		"x",
		"transform to apply (name@version)",
		func(value string) error {
			transformArgs = append(transformArgs, value)
			return nil
		},
	)

	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New(
			"usage: get [-p collection] [-o file] [-x name@version] <id>",
		)
	}

	transforms := make([]domain.Transform, 0, len(transformArgs))

	for _, value := range transformArgs {
		transform, err := parseTransform(value)
		if err != nil {
			return err
		}

		transforms = append(transforms, transform)
	}

	resp, err := port.Get(
		ctx,
		api.GetRequest{
			ID:         fs.Arg(0),
			Collection: *collection,
			Transforms: transforms,
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

func GetByHash(ctx context.Context, port app.ObjectPort, args []string) error {
	fs := flag.NewFlagSet("get-by-hash", flag.ContinueOnError)

	out := fs.String("o", "", "output file")
	collection := fs.String("c", "", "collection to retrieve object")

	var transformArgs []string

	fs.Func(
		"x",
		"transform to apply (name@version)",
		func(value string) error {
			transformArgs = append(transformArgs, value)
			return nil
		},
	)

	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New(
			"usage: get-by-hash [-p collection] [-o file] [-x name@version] <hash>",
		)
	}

	transforms := make([]domain.Transform, 0, len(transformArgs))

	for _, value := range transformArgs {
		transform, err := parseTransform(value)
		if err != nil {
			return err
		}

		transforms = append(transforms, transform)
	}

	resp, err := port.GetByHash(
		ctx,
		api.GetByHashRequest{
			Hash:       fs.Arg(0),
			Collection: *collection,
			Transforms: transforms,
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
