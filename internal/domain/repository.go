package domain

import "context"

type ObjectRepository interface {
    SaveObject(context.Context, *Object) error

	SavePayload(context.Context, *Payload) error
}

