package api

import (
	"MeshNet/internal/domain"
)

type PutRequest struct {
	Source     domain.Source
	Collection string
	Name       string
	MediaType  string
	Data       []byte
	Transforms []domain.Transform
}

type GetRequest struct {
	ID         string
	Collection string
	Transforms []domain.Transform
}

type GetByHashRequest struct {
	Hash       string
	Collection string
	Transforms []domain.Transform
}

type ListRequest struct {
	Collection string
}

type DeleteRequest struct {
	ID         string
	Collection string
}

type DeleteByHashRequest struct {
	Hash       string
	Collection string
}
