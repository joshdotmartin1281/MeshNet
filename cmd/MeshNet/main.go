package main

import (
	"fmt"
	"net/http"

	"MeshNet/internal/app"
	"MeshNet/internal/hash"
	"MeshNet/internal/processors/encrypt"
	"MeshNet/internal/processors/text"
	"MeshNet/internal/storage/sqlite"
	httptransport "MeshNet/internal/transport/http"
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

	httpHandler := httptransport.NewHandler(service)

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler.Routes(),
	}

	fmt.Println("HTTP server listening on :8080")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("HTTP server error:", err)
		return
	}
}
