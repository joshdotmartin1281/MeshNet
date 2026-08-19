package api

import (
	"MeshNet/internal/domain"
)

type PutRequest struct {
	Source     domain.Source
	Data       []byte
	Transforms []domain.Transform
}

type GetRequest struct {
	ID         string
	Transforms []domain.Transform
}

type GetByHashRequest struct {
	Hash       string
	Transforms []domain.Transform
}

type ListRequest struct{}

type DeleteRequest struct {
	ID string
}
