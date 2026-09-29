package api

import (
	"MeshNet/internal/domain"
	"MeshNet/transform"
)

type PutRequest struct {
	Source     domain.Source
	Collection string
	Name       string
	MediaType  string
	Data       []byte
	Transforms []transform.Transform
}

type GetRequest struct {
	ID         string
	Collection string
	Transforms []transform.Transform
}

type GetByHashRequest struct {
	Hash       string
	Collection string
	Transforms []transform.Transform
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
