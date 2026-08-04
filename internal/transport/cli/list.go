package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"MeshNet/internal/app"
)

func List(s *app.Service, args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	objects, err := s.List(context.Background())
	if err != nil {
		return err
	}

	for _, obj := range objects {
		fmt.Fprintf(os.Stdout, "%s\t%s\t%d\n",
			obj.ID,
			obj.Hash,
			obj.Size,
		)
	}

	return nil
}

