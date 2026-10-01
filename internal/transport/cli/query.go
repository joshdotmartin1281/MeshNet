package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/url"
	"os"
	"strings"

	"MeshNet/internal/app"
)

func parseParams(pairs []string) (url.Values, error) {
	params := url.Values{}

	for _, pair := range pairs {
		key, value, ok := strings.Cut(pair, "=")
		if !ok {
			return nil, errors.New(
				"invalid param " + pair + ": expected key=value",
			)
		}

		params.Add(key, value)
	}

	return params, nil
}

func Query(ctx context.Context, port app.RelationalPort, args []string) error {
	fs := flag.NewFlagSet("query", flag.ContinueOnError)

	var paramArgs []string

	fs.Func(
		"p",
		"query param (key=value)",
		func(value string) error {
			paramArgs = append(paramArgs, value)
			return nil
		},
	)

	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New(
			"usage: query [-p key=value]... <name>",
		)
	}

	params, err := parseParams(paramArgs)
	if err != nil {
		return err
	}

	resp, err := port.RunQuery(ctx, fs.Arg(0), params)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(resp)
}

func MutateQuery(ctx context.Context, port app.RelationalPort, args []string) error {
	fs := flag.NewFlagSet("mutate-query", flag.ContinueOnError)

	var paramArgs []string

	fs.Func(
		"p",
		"query param (key=value)",
		func(value string) error {
			paramArgs = append(paramArgs, value)
			return nil
		},
	)

	fs.SetOutput(io.Discard)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return errors.New(
			"usage: mutate-query [-p key=value]... <name>",
		)
	}

	params, err := parseParams(paramArgs)
	if err != nil {
		return err
	}

	resp, err := port.MutateQuery(ctx, fs.Arg(0), params)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(resp)
}
