package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"MeshNet/internal/app"
	"MeshNet/internal/storage/memory"
	"MeshNet/internal/transport/cli"
)

func main() {
	repo := memory.New()
	service := app.New(repo)

	files := []string{
		"test1.txt",
		"test2.txt",
		"test3.txt",
	}

	var ids []string
	var hashes []string

	for _, file := range files {
		out, err := captureOutput(func() error {
			return cli.Run(service, []string{"put", file})
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("=== PUT %s ===\n", file)
		fmt.Print(out)

		var id string
		var hash string

		fmt.Sscanf(out, "ID:   %s\nHash: %s", &id, &hash)

		ids = append(ids, id)
		hashes = append(hashes, hash)
	}

	fmt.Println("=== LIST ===")
	listOut, err := captureOutput(func() error {
		return cli.Run(service, []string{"list"})
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(listOut)

	fmt.Println("=== GET BY ID ===")
	for _, id := range ids {
		fmt.Printf("GET %s\n", id)

		if err := cli.Run(service, []string{"get", id}); err != nil {
			log.Fatal(err)
		}

		fmt.Println()
	}

	fmt.Println("=== GET BY HASH ===")
	for _, hash := range hashes {
		fmt.Printf("GET-BY-HASH %s\n", hash)

		if err := cli.Run(service, []string{"get-by-hash", hash}); err != nil {
			log.Fatal(err)
		}

		fmt.Println()
	}

	fmt.Println("=== DELETE ===")
	if err := cli.Run(service, []string{"delete", ids[1]}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== LIST AFTER DELETE ===")
	listOut, err = captureOutput(func() error {
		return cli.Run(service, []string{"list"})
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(listOut)

	fmt.Println("=== GET DELETED BY ID ===")
	if err := cli.Run(service, []string{"get", ids[1]}); err != nil {
		fmt.Println(err) // expect: object not found
	}

	fmt.Println("=== GET DELETED BY HASH ===")
	if err := cli.Run(service, []string{"get-by-hash", hashes[1]}); err != nil {
		fmt.Println(err) // expect: object not found
	}

	fmt.Println("=== VERIFY REMAINING ===")
	for _, id := range []string{ids[0], ids[2]} {
		fmt.Printf("GET %s\n", id)

		if err := cli.Run(service, []string{"get", id}); err != nil {
			log.Fatal(err)
		}

		fmt.Println()
	}

	fmt.Println("=== DELETE AGAIN ===")
	if err := cli.Run(service, []string{"delete", ids[1]}); err != nil {
		fmt.Println(err) // expect: object not found
	}
}

func captureOutput(fn func() error) (string, error) {
	old := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)

	return buf.String(), err
}
