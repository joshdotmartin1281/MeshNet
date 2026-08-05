package main

import (
	"context"
	"fmt"
	"log"

	"MeshNet/internal/app"
	"MeshNet/internal/domain"
	"MeshNet/internal/storage/file"
)

func main() {
	ctx := context.Background()

	repo, err := file.New("./meshnet-data")
	if err != nil {
		log.Fatal(err)
	}

	service := app.New(repo)

	data := []byte("hello MeshNet")

	created, err := service.Put(ctx, domain.SourceCLI, data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("created:")
	fmt.Println("id:", created.ID)
	fmt.Println("hash:", created.Hash)
	fmt.Println("source:", created.Source)
	fmt.Println("size:", created.Size)

	found, payload, err := service.Get(ctx, created.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nretrieved:")
	fmt.Println("id:", found.ID)
	fmt.Println("hash:", found.Hash)
	fmt.Println("payload:", string(payload.Data))
}

