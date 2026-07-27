package main

import (
	"context"
	"fmt"

	"MeshNet/internal/app"
	"MeshNet/internal/domain"
	"MeshNet/internal/storage/memory"
)

func main() {
	repo := memory.New()
	service := app.New(repo)

	obj := &domain.Object{
		ID:     "obj-1",
		Source: "cli",
	}

	payload := &domain.Payload{
		Data: []byte("hello world"),
	}

	if err := service.Process(context.Background(), payload, obj); err != nil {
		panic(err)
	}

	fmt.Printf("ID:         %s\n", obj.ID)
	fmt.Printf("Hash:       %s\n", obj.Hash)
	fmt.Printf("Size:       %d\n", obj.Size)
	fmt.Printf("Source:     %s\n", obj.Source)
	fmt.Printf("Created At: %s\n", obj.CreatedAt.Format("2006-01-02 15:04:05"))

	fmt.Printf("\nPayload ObjectID: %s\n", payload.ObjectID)
}

