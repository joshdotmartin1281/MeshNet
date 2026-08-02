package main

import (
	"log"
	"os"

	"MeshNet/internal/app"
	"MeshNet/internal/transport/cli"
	"MeshNet/internal/storage/memory"
)

func main() {
	repo := memory.New()
	service := app.New(repo)

	if err := cli.Run(service, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

