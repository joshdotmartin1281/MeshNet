package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"MeshNet/internal/app"
)

func Shell(ctx context.Context, port app.Port, in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	fmt.Fprintln(out, "MeshNet shell. Type 'exit' or 'quit' to leave.")

	for {
		fmt.Fprint(out, "> ")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if line == "exit" || line == "quit" {
			break
		}

		args := strings.Fields(line)

		if err := Run(ctx, port, args); err != nil {
			fmt.Fprintln(out, "error:", err)
		}
	}

	return scanner.Err()
}
