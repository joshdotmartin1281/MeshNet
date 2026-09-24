package http

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) runQuery(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	resp, err := h.port.RunQuery(r.Context(), name, r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) mutateQuery(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	resp, err := h.port.MutateQuery(r.Context(), name, r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
