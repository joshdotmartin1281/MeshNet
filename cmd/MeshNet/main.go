package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"MeshNet/internal/app"
	"MeshNet/internal/storage/memory"
	"MeshNet/internal/transport/cli"
)

func main() {
	repo := memory.New()
	service := app.New(repo)

	var buf bytes.Buffer

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cli.Run(service, []string{"put", "test.txt"})
	if err != nil {
		log.Fatal(err)
	}

	w.Close()
	os.Stdout = old

	buf.ReadFrom(r)

	fmt.Print("=== PUT ===\n")
	fmt.Print(buf.String())

	// parse ID from output
	var id string
	fmt.Sscanf(buf.String(), "ID: %s", &id)

	fmt.Println("=== GET ===")

	if err := cli.Run(service, []string{"get", id}); err != nil {
		log.Fatal(err)
	}
}

