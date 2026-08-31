package main

import (
	"context"
	"fmt"

	"MeshNet/internal/api"
	"MeshNet/internal/app"
	"MeshNet/internal/domain"
	"MeshNet/internal/hash"
	"MeshNet/internal/processors/encrypt"
	"MeshNet/internal/processors/text"
	"MeshNet/internal/storage/file"
)

func main() {
	repo, err := file.New("./mesh-net")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	hasher := hash.NewSHA256()

	key, err := encrypt.GenerateKey()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	processor := app.NewProcessor(
		text.NewUppercase(),
		text.NewLowercase(),
		encrypt.NewEncrypt(key),
		encrypt.NewDecrypt(key),
	)

	service := app.New(repo, hasher, processor)

	ctx := context.Background()

	putAlice, err := service.Put(ctx, api.PutRequest{
		Path:   "alice",
		Source: domain.SourceMANUAL,
		Data:   []byte("hello from alice"),
	})
	if err != nil {
		fmt.Println("put alice error:", err)
		return
	}

	fmt.Println("Alice object:")
	fmt.Printf("ID:   %s\n", putAlice.Object.ID)
	fmt.Printf("Path: %s\n", putAlice.Object.Path)
	fmt.Printf("Hash: %s\n", putAlice.Object.Hash)
	fmt.Printf("Size: %d\n", putAlice.Object.Size)
	fmt.Println()

	putBob, err := service.Put(ctx, api.PutRequest{
		Path:   "bob",
		Source: domain.SourceMANUAL,
		Data:   []byte("hello from bob"),
	})
	if err != nil {
		fmt.Println("put bob error:", err)
		return
	}

	fmt.Println("Bob object:")
	fmt.Printf("ID:   %s\n", putBob.Object.ID)
	fmt.Printf("Path: %s\n", putBob.Object.Path)
	fmt.Printf("Hash: %s\n", putBob.Object.Hash)
	fmt.Printf("Size: %d\n", putBob.Object.Size)
	fmt.Println()

	alice, err := service.Get(ctx, api.GetRequest{
		Path: "alice",
		ID:   putAlice.Object.ID,
	})
	if err != nil {
		fmt.Println("get alice error:", err)
		return
	}

	fmt.Println("Retrieved Alice object:")
	fmt.Printf("ID:   %s\n", alice.Object.ID)
	fmt.Printf("Path: %s\n", alice.Object.Path)
	fmt.Printf("Data: %s\n", alice.Payload.Data)
	fmt.Println()

	aliceObjects, err := service.List(ctx, api.ListRequest{
		Path: "alice",
	})
	if err != nil {
		fmt.Println("list alice error:", err)
		return
	}

	fmt.Println("Alice objects:")

	for _, obj := range aliceObjects.Objects {
		fmt.Printf(
			"ID: %s | Path: %s | Hash: %s\n",
			obj.ID,
			obj.Path,
			obj.Hash,
		)
	}

	fmt.Println()

	bobObjects, err := service.List(ctx, api.ListRequest{
		Path: "bob",
	})
	if err != nil {
		fmt.Println("list bob error:", err)
		return
	}

	fmt.Println("Bob objects:")

	for _, obj := range bobObjects.Objects {
		fmt.Printf(
			"ID: %s | Path: %s | Hash: %s\n",
			obj.ID,
			obj.Path,
			obj.Hash,
		)
	}
}
