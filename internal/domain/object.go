package domain

import (
	"time"

	"MeshNet/transform"
)

type Object struct {
	ID         string
	Hash       string
	Collection string
	Name       string
	MediaType  string
	Size       int64
	Source     Source
	CreatedAt  time.Time
	Transforms []transform.Transform
}

type Payload struct {
	ObjectID string
	Data     []byte
}
