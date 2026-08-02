
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"MeshNet/internal/app"
	"MeshNet/internal/domain"
)

func Put(s *app.Service, args []string) error {
	fs := flag.NewFlagSet("put", flag.ContinueOnError)

	file := fs.String("f", "", "file to upload")
	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	var (
		data []byte
		err error
	)

	switch {
	case *file != "":
		data, err = os.ReadFile(*file)
	case fs.NArg() == 1:
		data, err = os.ReadFile(fs.Arg(0))
	default:
		data, err = io.ReadAll(os.Stdin)
	}

	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("no input data")
	}

	obj, err := s.Process(context.Background(), domain.SourceCLI, data)
	if err != nil {
		return err
	}

	fmt.Printf("ID:   %s\n", obj.ID)
	fmt.Printf("Hash: %s\n", obj.Hash)
	fmt.Printf("Size: %d bytes\n", obj.Size)

	return nil
}

