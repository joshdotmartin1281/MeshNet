package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"MeshNet/internal/api"
	"MeshNet/internal/app"
	"MeshNet/internal/domain"
)

func Put(port app.Port, args []string) error {
	fs := flag.NewFlagSet("put", flag.ContinueOnError)

	file := fs.String("f", "", "file to upload")

	var transformArgs []string
	fs.Func(
		"transform",
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

	var (
		data []byte
		err  error
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

	transforms := make([]domain.Transform, 0, len(transformArgs))

	for _, value := range transformArgs {
		transform, err := parseTransform(value)
		if err != nil {
			return err
		}

		transforms = append(transforms, transform)
	}

	resp, err := port.Put(
		context.Background(),
		api.PutRequest{
			Source:     domain.SourceCLI,
			Data:       data,
			Transforms: transforms,
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("ID:   %s\n", resp.Object.ID)
	fmt.Printf("Hash: %s\n", resp.Object.Hash)
	fmt.Printf("Size: %d bytes\n", resp.Object.Size)

	return nil
}

func parseTransform(value string) (domain.Transform, error) {
	parts := strings.SplitN(value, "@", 2)

	if len(parts) != 2 {
		return domain.Transform{}, fmt.Errorf(
			"invalid transform %q: expected name@version",
			value,
		)
	}

	if parts[0] == "" || parts[1] == "" {
		return domain.Transform{}, fmt.Errorf(
			"invalid transform %q: expected name@version",
			value,
		)
	}

	return domain.Transform{
		Name:    parts[0],
		Version: parts[1],
		Params:  map[string]string{},
	}, nil
}
