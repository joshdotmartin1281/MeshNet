package http

import (
	"net/http"
	"strings"

	"MeshNet/uauth"
)

func (h *Handler) protected(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !h.port.AuthEnabled() {
			next(w, r)
			return
		}

		token, ok := bearerToken(r)
		if !ok {
			writeError(w, auth.ErrUnauthorized)
			return
		}

		id, err := h.port.Authenticate(r.Context(), token)
		if err != nil {
			writeError(w, err)
			return
		}

		next(w, r.WithContext(auth.WithIdentity(r.Context(), id)))
	})
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	value := r.Header.Get("Authorization")

	if len(value) <= len(prefix) || !strings.EqualFold(value[:len(prefix)], prefix) {
		return "", false
	}

	return strings.TrimSpace(value[len(prefix):]), true
}
