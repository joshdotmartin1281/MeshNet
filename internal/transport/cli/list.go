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

func List(port app.Port, args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	resp, err := port.List(
		context.Background(),
		api.ListRequest{},
	)
	if err != nil {
		return err
	}

	for _, obj := range resp.Objects {
		fmt.Fprintf(os.Stdout, "%s\t%s\t%d\n",
			obj.ID,
			obj.Hash,
			obj.Size,
		)
	}

	return nil
}

