package main

import (
	"context"
	"fmt"
	"log"

	"MeshNet/internal/app"
	"MeshNet/internal/domain"
	"MeshNet/internal/hash"
	"MeshNet/internal/storage/memory"
)

func main() {
	ctx := context.Background()

	repo := memory.New()

	hasher := hash.NewSHA256()

	service := app.New(repo, hasher)

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

