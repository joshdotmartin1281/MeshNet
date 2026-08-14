package api

import "MeshNet/internal/domain"

type PutResponse struct {
	Object *domain.Object
}

type GetResponse struct {
	Object  *domain.Object
	Payload *domain.Payload
}

type ListResponse struct {
	Objects []*domain.Object
}

type DeleteResponse struct{}

