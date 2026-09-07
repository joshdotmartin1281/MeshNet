package domain

import "errors"

var (
	ErrNotFound  = errors.New("object not found")
	ErrDuplicate = errors.New("object already exists")
)
