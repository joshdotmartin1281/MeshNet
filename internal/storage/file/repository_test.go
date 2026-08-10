package file_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"MeshNet/internal/domain"
	"MeshNet/internal/storage/file"
)

func TestStore_PersistsAcrossInstances(t *testing.T) {
	dir := t.TempDir()

	ctx := context.Background()

	store1, err := file.New(dir)
	if err != nil {
		t.Fatalf("file.New() error = %v", err)
	}

	obj := &domain.Object{
		ID:   "object-1",
		Hash: "hash-1",
	}

	payload := &domain.Payload{
		ObjectID: obj.ID,
		Data:     []byte("hello"),
	}

	if err := store1.Save(
		ctx,
		obj,
		payload,
	); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// the same directory.
	store2, err := file.New(dir)
	if err != nil {
		t.Fatalf("file.New() error = %v", err)
	}

	gotObject, gotPayload, err := store2.Get(
		ctx,
		obj.ID,
	)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if gotObject.ID != obj.ID {
		t.Errorf(
			"ID = %q, want %q",
			gotObject.ID,
			obj.ID,
		)
	}

	if string(gotPayload.Data) != string(payload.Data) {
		t.Errorf(
			"Data = %q, want %q",
			gotPayload.Data,
			payload.Data,
		)
	}
}

func TestStore_CreatesDirectories(t *testing.T) {
	dir := t.TempDir()

	storeDir := filepath.Join(
		dir,
		"meshnet",
	)

	_, err := file.New(storeDir)
	if err != nil {
		t.Fatalf("file.New() error = %v", err)
	}

	objectsDir := filepath.Join(
		storeDir,
		"objects",
	)

	hashesDir := filepath.Join(
		storeDir,
		"hashes",
	)

	if info, err := os.Stat(objectsDir); err != nil {
		t.Fatalf(
			"objects directory does not exist: %v",
			err,
		)
	} else if !info.IsDir() {
		t.Fatalf(
			"objects path is not a directory",
		)
	}

	if info, err := os.Stat(hashesDir); err != nil {
		t.Fatalf(
			"hashes directory does not exist: %v",
			err,
		)
	} else if !info.IsDir() {
		t.Fatalf(
			"hashes path is not a directory",
		)
	}
}

func TestStore_WritesObjectAndHashFiles(t *testing.T) {
	dir := t.TempDir()

	store, err := file.New(dir)
	if err != nil {
		t.Fatalf("file.New() error = %v", err)
	}

	obj := &domain.Object{
		ID:   "object-1",
		Hash: "hash-1",
	}

	payload := &domain.Payload{
		ObjectID: obj.ID,
		Data:     []byte("hello"),
	}

	if err := store.Save(
		context.Background(),
		obj,
		payload,
	); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	objectPath := filepath.Join(
		dir,
		"objects",
		obj.ID+".json",
	)

	hashPath := filepath.Join(
		dir,
		"hashes",
		obj.Hash,
	)

	if _, err := os.Stat(objectPath); err != nil {
		t.Fatalf(
			"object file does not exist: %v",
			err,
		)
	}

	if _, err := os.Stat(hashPath); err != nil {
		t.Fatalf(
			"hash file does not exist: %v",
			err,
		)
	}
}

func TestStore_DeleteRemovesFiles(t *testing.T) {
	dir := t.TempDir()

	store, err := file.New(dir)
	if err != nil {
		t.Fatalf("file.New() error = %v", err)
	}

	obj := &domain.Object{
		ID:   "object-1",
		Hash: "hash-1",
	}

	payload := &domain.Payload{
		ObjectID: obj.ID,
		Data:     []byte("hello"),
	}

	ctx := context.Background()

	if err := store.Save(
		ctx,
		obj,
		payload,
	); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if err := store.Delete(ctx, obj.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	objectPath := filepath.Join(
		dir,
		"objects",
		obj.ID+".json",
	)

	hashPath := filepath.Join(
		dir,
		"hashes",
		obj.Hash,
	)

	if _, err := os.Stat(objectPath); !os.IsNotExist(err) {
		t.Fatalf(
			"object file still exists, err = %v",
			err,
		)
	}

	if _, err := os.Stat(hashPath); !os.IsNotExist(err) {
		t.Fatalf(
			"hash file still exists, err = %v",
			err,
		)
	}
}

