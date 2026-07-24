package main

import (
	"context"

	"MeshNet/internal/app"
	"MeshNet/internal/domain"
	"MeshNet/internal/storage/memory"
)

func main() {
	repo := memory.New()

	service := app.New(repo)

	msg := &domain.Message{
		Source: "cli",
		Data:   []byte("hello world"),
	}

	service.Process(context.Background(), msg)
}

