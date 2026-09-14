package http

import (
	"net/http"

	"MeshNet/internal/api"
)

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	collection := r.URL.Query().Get("collection")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	resp, err := h.port.Get(
		r.Context(),
		api.GetRequest{
			ID:         id,
			Collection: collection,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(resp.Payload.Data)
}
