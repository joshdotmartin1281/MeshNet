package domain

import "time"

type Transform struct {
	Name    string
	Version string
	Value   string
}

type Message struct {
	ID          string
	Source      string
	Hash        string
	Path        string
	Size		int64
	Transforms  []Transform
	CreatedAt   time.Time
}

type Payload struct {
	Data []byte
}

type MessageResponse struct {
	ID        string
	Size      int64
	CreatedAt time.Time
}

