package domain

import "github.com/google/uuid"

type ID struct {
	uuid.UUID
}

func NewID() ID {
	return ID{UUID: uuid.New()}
}
