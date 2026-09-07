package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"MeshNet/internal/api"
	"MeshNet/internal/app"
	"MeshNet/internal/domain"
)

func Put(ctx context.Context, port app.ObjectPort, args []string) error {
	fs := flag.NewFlagSet("put", flag.ContinueOnError)

	file := fs.String("f", "", "file to upload")
	text := fs.String("t", "", "text to input")
	collection := fs.String("c", "", "collection to store object")

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

	if strings.TrimSpace(*collection) == "" {
		return errors.New("collection is required")
	}

	if *file != "" && *text != "" {
		return errors.New("cannot use -f and -t together")
	}

	if fs.NArg() > 0 {
		return fmt.Errorf(
			"unexpected argument %q: use -f for a file or -t for text",
			fs.Arg(0),
		)
	}

	var (
		data      []byte
		err       error
		name      string
		mediaType string
	)

	switch {
	case *file != "":
		data, err = os.ReadFile(*file)
		if err != nil {
			return err
		}

		name = filepath.Base(*file)

	case *text != "":
		data = []byte(*text)
		name = "text.txt"

	default:
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		name = "stdin"
	}

	if len(data) == 0 {
		return errors.New("no input data")
	}

	if *text != "" {
		mediaType = "text/plain"
	} else {
		mediaType = http.DetectContentType(data)
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
		ctx,
		api.PutRequest{
			Source:     domain.SourceCLI,
			Collection: *collection,
			Name:       name,
			MediaType:  mediaType,
			Data:       data,
			Transforms: transforms,
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("ID:         %s\n", resp.Object.ID)
	fmt.Printf("Collection: %s\n", resp.Object.Collection)
	fmt.Printf("Name:       %s\n", resp.Object.Name)
	fmt.Printf("MediaType:  %s\n", resp.Object.MediaType)
	fmt.Printf("Hash:       %s\n", resp.Object.Hash)
	fmt.Printf("Size:       %d bytes\n", resp.Object.Size)

	return nil
}

func parseTransform(value string) (domain.Transform, error) {
	value = strings.TrimSpace(value)

	parts := strings.SplitN(value, "@", 2)

	if len(parts) != 2 {
		return domain.Transform{}, fmt.Errorf(
			"invalid transform %q: expected name@version",
			value,
		)
	}

	name := strings.TrimSpace(parts[0])
	version := strings.TrimSpace(parts[1])

	if name == "" || version == "" {
		return domain.Transform{}, fmt.Errorf(
			"invalid transform %q: expected name@version",
			value,
		)
	}

	return domain.Transform{
		Name:    name,
		Version: version,
		Params:  map[string]string{},
	}, nil
}
