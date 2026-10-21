package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"MeshNet/cmd/example0/tags"
	"MeshNet/internal/app"
	"MeshNet/internal/hash"
	"MeshNet/internal/storage/sqlite"
	"MeshNet/internal/transport/cli"
	httptransport "MeshNet/internal/transport/http"
	"MeshNet/processors/encrypt"
	"MeshNet/processors/text"
	"MeshNet/uauth"
	authfile "MeshNet/uauth/file"
)

func main() {
	repo, err := sqlite.New("./cmd/example0/sqlite")
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
	authStore, err := authfile.New("./auth/users.json")
	if err != nil {
		fmt.Println("auth init error:", err)
		return
	}

	service := app.New(repo, hasher, processor, queries, app.WithAuth(authStore))

	httpHandler := httptransport.NewHandler(service)

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler.Routes("./cmd/example0/web"),
	}

	go func() {
		fmt.Println("HTTP server listening on :8080")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP server error:", err)
		}
	}()
	operatorCtx := auth.WithIdentity(context.Background(), auth.Operator)

	if err := cli.Shell(operatorCtx, service, os.Stdin, os.Stdout); err != nil {
		fmt.Println("shell error:", err)
	}

	_ = server.Shutdown(context.Background())
}
