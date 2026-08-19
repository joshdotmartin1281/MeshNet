package cli

import "errors"

var (
	ErrNoCommand      = errors.New("no command specified")
	ErrUnknownCommand = errors.New("unknown command")
)
