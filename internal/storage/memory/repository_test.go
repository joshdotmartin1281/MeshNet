package memory_test

import (
	"testing"

	"MeshNet/internal/storage/memory"
)

func TestNew(t *testing.T) {
	store := memory.New()

	if store == nil {
		t.Fatal("New() returned nil")
	}
}

