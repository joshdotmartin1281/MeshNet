package http

import (
	"encoding/json"
	"io"
	"net/http"

	"MeshNet/internal/api"
	"MeshNet/internal/domain"
)

func (h *Handler) put(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(data) == 0 {
		http.Error(w, "no input data", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()

	collection := query.Get("collection")
	if collection == "" {
		http.Error(w, "collection is required", http.StatusBadRequest)
		return
	}

	name := query.Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	mediaType := r.Header.Get("Content-Type")
	if mediaType == "" {
		mediaType = http.DetectContentType(data)
	}

	transforms := make([]domain.Transform, 0, len(query["transform"]))

	for _, value := range query["transform"] {
		transform, err := domain.ParseTransform(value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		transforms = append(transforms, transform)
	}

	resp, err := h.port.Put(
		r.Context(),
		api.PutRequest{
			Source:     domain.SourceHTTP,
			Collection: collection,
			Name:       name,
			MediaType:  mediaType,
			Data:       data,
			Transforms: transforms,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(resp)
}
