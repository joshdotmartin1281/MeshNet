package domain

import "context"

type Repository interface {
	Save(context.Context, *Message) error
}

