package http

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /objects", h.put)

	mux.HandleFunc("GET /objects/{id}", h.get)
	mux.HandleFunc("GET /objects/hash/{hash}", h.getByHash)
	mux.HandleFunc("GET /objects", h.list)

	mux.HandleFunc("DELETE /objects/{id}", h.delete)
	mux.HandleFunc("DELETE /objects/hash/{hash}", h.deleteByHash)

	return mux
}
