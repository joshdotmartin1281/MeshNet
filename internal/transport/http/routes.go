package http

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	//mux.HandleFunc("PUT /objects/{id}", h.put)
	//mux.HandleFunc("GET /objects/{id}", h.get)
	//mux.HandleFunc("GET /objects", h.list)
	//mux.HandleFunc("DELETE /objects/{id}", h.delete)

	return mux
}
