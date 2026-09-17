package http

import (
	"errors"
	"net/http"

	"MeshNet/internal/domain"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, domain.ErrDuplicate):
		http.Error(w, err.Error(), http.StatusConflict)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
