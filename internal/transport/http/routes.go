package http

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /objects", h.put)
	mux.HandleFunc("GET /objects/{id}", h.get)

	return mux
}
