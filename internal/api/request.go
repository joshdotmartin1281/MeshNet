package api

import (
	"MeshNet/internal/domain"
)

type PutRequest struct {
	Source     domain.Source
	Path	   string
	Data       []byte
	Transforms []domain.Transform
}

type GetRequest struct {
	ID         string
	Path	   string
	Transforms []domain.Transform
}

type GetByHashRequest struct {
	Hash       string
	Path       string
	Transforms []domain.Transform
}

type ListRequest struct{
	Path	   string
}

type DeleteRequest struct {
	ID 		string
	Path 	string
}

type DeleteByHashRequest struct {
	Hash 	string
	Path 	string
}
