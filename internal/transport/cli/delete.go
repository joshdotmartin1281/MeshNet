package cli

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"

	"MeshNet/internal/api"
	"MeshNet/internal/app"
)

func Delete(ctx context.Context, port app.Port, args []string) error {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)

	collection := fs.String(
		"c",
		"",
		"collection containing object",
	)

	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*collection) == "" {
		return errors.New("collection is required")
	}

	if fs.NArg() != 1 {
		return errors.New("usage: delete -c <collection> <id>")
	}

	_, err := port.Delete(
		ctx,
		api.DeleteRequest{
			ID:         fs.Arg(0),
			Collection: *collection,
		},
	)

	return err
}

func DeleteByHash(ctx context.Context, port app.Port, args []string) error {
	fs := flag.NewFlagSet("delete-by-hash", flag.ContinueOnError)

	collection := fs.String(
		"c",
		"",
		"collection containing object",
	)

	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*collection) == "" {
		return errors.New("collection is required")
	}

	if fs.NArg() != 1 {
		return errors.New(
			"usage: delete-by-hash -c <collection> <hash>",
		)
	}

	_, err := port.DeleteByHash(
		ctx,
		api.DeleteByHashRequest{
			Hash:       fs.Arg(0),
			Collection: *collection,
		},
	)

	return err
}
