package http

import "MeshNet/internal/app"

type Handler struct {
	port app.Port
}

func NewHandler(port app.Port) *Handler {
	return &Handler {
		port: port,
	}
}
