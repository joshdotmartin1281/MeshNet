package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "strings"

    "MeshNet/internal/app"
    "MeshNet/internal/domain"
    "MeshNet/internal/storage/memory"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatalf("usage: meshctl <text>")
    }

    data := []byte(strings.Join(os.Args[1:], " "))

    repo := memory.New()
    service := app.New(repo)

    obj, err := service.Process(
        context.Background(),
        domain.SourceCLI,
        data,
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %s\n", obj.ID)
    fmt.Printf("Hash: %s\n", obj.Hash)
}

