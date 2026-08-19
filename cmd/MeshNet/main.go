package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"MeshNet/internal/app"
	"MeshNet/internal/hash"
	"MeshNet/internal/processors/text"
	"MeshNet/internal/storage/file"
	"MeshNet/internal/transport/cli"
)

func main() {
	repo, err := file.New("./mesh-net")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	hasher := hash.NewSHA256()

	processor := app.NewProcessor(
		text.NewUppercase(),
		text.NewLowercase(),
	)

	service := app.New(repo, hasher, processor)

	fmt.Println("MeshNet")
	fmt.Println("Type a command, or 'exit' to quit.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

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

		if err := cli.Run(service, args); err != nil {
			fmt.Println("error:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error:", err)
	}
}
