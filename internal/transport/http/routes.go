package http

import "net/http"

func (h *Handler) Routes(path string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /objects", h.put)

	mux.HandleFunc("GET /objects/{id}", h.get)
	mux.HandleFunc("GET /objects/hash/{hash}", h.getByHash)
	mux.HandleFunc("GET /objects", h.list)
	mux.HandleFunc("GET /health", h.health)

	mux.HandleFunc("DELETE /objects/{id}", h.delete)
	mux.HandleFunc("DELETE /objects/hash/{hash}", h.deleteByHash)

	mux.Handle("/", http.FileServer(http.Dir(path)))	

	return mux
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
