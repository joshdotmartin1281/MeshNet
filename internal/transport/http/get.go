package http

import (
	"net/http"

	"MeshNet/internal/api"
	"MeshNet/internal/domain"
)

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	collection := r.URL.Query().Get("collection")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	transforms := make([]domain.Transform, 0, len(r.URL.Query()["transform"]))

	for _, value := range r.URL.Query()["transform"] {
		transform, err := domain.ParseTransform(value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		transforms = append(transforms, transform)
	}

	resp, err := h.port.Get(
		r.Context(),
		api.GetRequest{
			ID:         id,
			Collection: collection,
			Transforms: transforms,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writePayload(w, resp)
}

func (h *Handler) getByHash(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	collection := r.URL.Query().Get("collection")

	if hash == "" {
		http.Error(w, "hash is required", http.StatusBadRequest)
		return
	}

	transforms := make([]domain.Transform, 0, len(r.URL.Query()["transform"]))

	for _, value := range r.URL.Query()["transform"] {
		transform, err := domain.ParseTransform(value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		transforms = append(transforms, transform)
	}

	resp, err := h.port.GetByHash(
		r.Context(),
		api.GetByHashRequest{
			Hash:       hash,
			Collection: collection,
			Transforms: transforms,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writePayload(w, resp)
}

func writePayload(w http.ResponseWriter, resp api.GetResponse) {
	if resp.Payload == nil {
		http.Error(w, "payload not found", http.StatusNotFound)
		return
	}

	if resp.Object != nil && resp.Object.MediaType != "" {
		w.Header().Set("Content-Type", resp.Object.MediaType)
	}

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(resp.Payload.Data)
}
