package main

import (
	"context"
	"fmt"
	"net/http"

	"MeshNet/internal/app"
	"MeshNet/internal/hash"
	"MeshNet/internal/processors/encrypt"
	"MeshNet/internal/processors/text"
	"MeshNet/internal/storage/sqlite"
	httptransport "MeshNet/internal/transport/http"
	"MeshNet/queries/tags"
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

	queries, err := app.NewQueries(context.Background(), repo,
		tags.NewSearchByTag(),
	)
	if err != nil {
		fmt.Println("query init error:", err)
		return
	}

	service := app.New(repo, hasher, processor, queries)

	httpHandler := httptransport.NewHandler(service)

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler.Routes("./web"),
	}

	fmt.Println("HTTP server listening on :8080")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("HTTP server error:", err)
		return
	}
}
