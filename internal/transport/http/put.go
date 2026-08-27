package http

import (
	"net/http"

	"MeshNet/internal/api"
	"MeshNet/internal/domain"
)

func (h *Handler) put(w http.ResponseWriter, r *http.Request) {
	// Decode HTTP request

	// Build application request
	req := api.PutRequest{
		Source:     domain.SourceHTTP,
		Data:       data,
		Transforms: transforms,
	}

	// Call application layer
	resp, err := h.port.Put(r.Context(), req)
	if err != nil {
		// HTTP error response
		return
	}

	// Encode response
}
