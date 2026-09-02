package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"MeshNet/internal/api"
	"MeshNet/internal/app"
)

func List(ctx context.Context, port app.Port, args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)

	collection := fs.String(
		"c",
		"",
		"collection to list",
	)

	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() > 0 {
		return fmt.Errorf(
			"unexpected argument %q",
			fs.Arg(0),
		)
	}

	resp, err := port.List(
		ctx,
		api.ListRequest{
			Collection: *collection,
		},
	)
	if err != nil {
		return err
	}

	for _, obj := range resp.Objects {
		fmt.Fprintf(
			os.Stdout,
			"%s\t%s\t%s\t%s\t%d\t%s\n",
			obj.ID,
			obj.Collection,
			obj.Name,
			obj.MediaType,
			obj.Size,
			obj.Hash,
		)
	}

	return nil
}
