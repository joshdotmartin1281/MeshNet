package domain

import "time"

type Object struct {
	ID         string
	Hash       string
	Size       int64
	Source     Source
	CreatedAt  time.Time
	Transforms []Transform
}

type Payload struct {
	ObjectID string
	Data     []byte
}
