package http

import (
	"net/http"

	"MeshNet/internal/api"
)

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	collection := r.URL.Query().Get("collection")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	_, err := h.port.Delete(
		r.Context(),
		api.DeleteRequest{
			ID:         id,
			Collection: collection,
		},
	)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteByHash(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	collection := r.URL.Query().Get("collection")

	if hash == "" {
		http.Error(w, "hash is required", http.StatusBadRequest)
		return
	}

	_, err := h.port.DeleteByHash(
		r.Context(),
		api.DeleteByHashRequest{
			Hash:       hash,
			Collection: collection,
		},
	)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
