package api

import (
	"MeshNet/internal/domain"
)

type PutRequest struct {
	Source domain.Source
	Data   []byte
}

type GetRequest struct {
	ID string
}

type GetByHashRequest struct {
	Hash string
}

type ListRequest struct{}

type DeleteRequest struct {
	ID string
}

