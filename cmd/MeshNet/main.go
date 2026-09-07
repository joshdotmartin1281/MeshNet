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
	"MeshNet/internal/storage"
	"MeshNet/internal/storage/sqlite"
	"MeshNet/internal/transport/cli"
)

func main() {
	repo, err := sqlite.New("./sqlite/mesh-net")
	if err != nil {
		fmt.Println("storage error:", err)
		return
	}
	defer repo.Close()

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
	var relational storage.RelationalStore = repo

	_, err = relational.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	)
`)
	if err != nil {
		fmt.Println("relational create table error:", err)
		return
	}

	_, err = relational.ExecContext(ctx, `
	INSERT INTO tags (name) VALUES (?)
`, "test")
	if err != nil {
		fmt.Println("relational insert error:", err)
		return
	}

	var tag string

	err = relational.QueryRowContext(ctx, `
	SELECT name
	FROM tags
	WHERE id = 1
`).Scan(&tag)
	if err != nil {
		fmt.Println("relational query error:", err)
		return
	}

	fmt.Println("Relational store:")
	fmt.Printf("Tag: %s\n", tag)
	fmt.Println()

	putAlice, err := service.Put(ctx, api.PutRequest{
		Collection: "alice",
		Source:     domain.SourceMANUAL,
		Name:       "alice.txt",
		MediaType:  "text/plain",
		Data:       []byte("hello from alice"),
	})
	if err != nil {
		fmt.Println("put alice error:", err)
		return
	}

	fmt.Println("Alice object:")
	fmt.Printf("ID:         %s\n", putAlice.Object.ID)
	fmt.Printf("Collection: %s\n", putAlice.Object.Collection)
	fmt.Printf("Name:       %s\n", putAlice.Object.Name)
	fmt.Printf("MediaType:  %s\n", putAlice.Object.MediaType)
	fmt.Printf("Hash:       %s\n", putAlice.Object.Hash)
	fmt.Printf("Size:       %d\n", putAlice.Object.Size)
	fmt.Println()

	putBob, err := service.Put(ctx, api.PutRequest{
		Collection: "bob",
		Source:     domain.SourceMANUAL,
		Name:       "bob.txt",
		MediaType:  "text/plain",
		Data:       []byte("hello from bob"),
	})
	if err != nil {
		fmt.Println("put bob error:", err)
		return
	}

	fmt.Println("Bob object:")
	fmt.Printf("ID:         %s\n", putBob.Object.ID)
	fmt.Printf("Collection: %s\n", putBob.Object.Collection)
	fmt.Printf("Name:       %s\n", putBob.Object.Name)
	fmt.Printf("MediaType:  %s\n", putBob.Object.MediaType)
	fmt.Printf("Hash:       %s\n", putBob.Object.Hash)
	fmt.Printf("Size:       %d\n", putBob.Object.Size)
	fmt.Println()

	alice, err := service.Get(ctx, api.GetRequest{
		Collection: "alice",
		ID:         putAlice.Object.ID,
	})
	if err != nil {
		fmt.Println("get alice error:", err)
		return
	}

	fmt.Println("Retrieved Alice object:")
	fmt.Printf("ID:         %s\n", alice.Object.ID)
	fmt.Printf("Collection: %s\n", alice.Object.Collection)
	fmt.Printf("Name:       %s\n", alice.Object.Name)
	fmt.Printf("MediaType:  %s\n", alice.Object.MediaType)
	fmt.Printf("Data:       %s\n", alice.Payload.Data)
	fmt.Println()

	aliceObjects, err := service.List(ctx, api.ListRequest{
		Collection: "alice",
	})
	if err != nil {
		fmt.Println("list alice error:", err)
		return
	}

	fmt.Println("Alice objects:")

	for _, obj := range aliceObjects.Objects {
		fmt.Printf(
			"ID: %s | Collection: %s | Name: %s | MediaType: %s | Hash: %s\n",
			obj.ID,
			obj.Collection,
			obj.Name,
			obj.MediaType,
			obj.Hash,
		)
	}

	fmt.Println()

	bobObjects, err := service.List(ctx, api.ListRequest{
		Collection: "bob",
	})
	if err != nil {
		fmt.Println("list bob error:", err)
		return
	}

	fmt.Println("Bob objects:")

	for _, obj := range bobObjects.Objects {
		fmt.Printf(
			"ID: %s | Collection: %s | Name: %s | MediaType: %s | Hash: %s\n",
			obj.ID,
			obj.Collection,
			obj.Name,
			obj.MediaType,
			obj.Hash,
		)
	}

	fmt.Println()

	fmt.Println("CLI Get:")

	if err := cli.Get(ctx, service, []string{
		"-c", "alice",
		putAlice.Object.ID,
	}); err != nil {
		fmt.Println("cli get error:", err)
		return
	}

	fmt.Println()

	fmt.Println("CLI List:")

	if err := cli.List(ctx, service, []string{
		"-c", "alice",
	}); err != nil {
		fmt.Println("cli list error:", err)
		return
	}

	fmt.Println()

	fmt.Println("CLI Delete:")

	if err := cli.Delete(ctx, service, []string{
		"-c", "alice",
		putAlice.Object.ID,
	}); err != nil {
		fmt.Println("cli delete error:", err)
		return
	}

	fmt.Println("Alice deleted.")
}
