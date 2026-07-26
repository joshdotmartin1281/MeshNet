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

	obj := &domain.Object{
		Source: "cli",
	}

	payload := &domain.Payload{
		Data: []byte("hello world"),
	}

	err := service.Process(context.Background(), payload, obj)
	if err != nil {
		panic(err)
	}
}

