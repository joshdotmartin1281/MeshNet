package http

import (
	"errors"
	"net/http"

	"MeshNet/internal/domain"
	"MeshNet/uauth"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, domain.ErrDuplicate):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, auth.ErrUnauthorized):
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, "unauthorized", http.StatusUnauthorized)

	case errors.Is(err, auth.ErrForbidden):
		http.Error(w, "forbidden", http.StatusForbidden)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
