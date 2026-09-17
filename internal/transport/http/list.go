package http

import (
	"encoding/json"
	"net/http"

	"MeshNet/internal/api"
)

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	collection := r.URL.Query().Get("collection")

	resp, err := h.port.List(
		r.Context(),
		api.ListRequest{
			Collection: collection,
		},
	)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}
