package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotSupported = errors.New("auth is not enabled or does not support user management")
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type Level int

const (
	Read  Level = 1
	Write Level = 2
)

func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "read":
		return Read, nil
	case "write":
		return Write, nil
	default:
		return 0, fmt.Errorf("invalid level %q: expected read or write", s)
	}
}

func (l Level) String() string {
	switch l {
	case Read:
		return "read"
	case Write:
		return "write"
	default:
		return "none"
	}
}

type Grant struct {
	Collection string
	Level      Level
}

type User struct {
	Name      string
	CreatedAt time.Time
	Grants    []Grant
}

type Identity struct {
	UserID string
	Admin  bool
}

var Operator = Identity{UserID: "operator", Admin: true}

type identityKey struct{}

func WithIdentity(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, identityKey{}, id)
}

func FromContext(ctx context.Context) (Identity, bool) {
	id, ok := ctx.Value(identityKey{}).(Identity)
	return id, ok
}

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (Identity, error)
}

type Authorizer interface {
	Authorize(ctx context.Context, collection string, need Level) error
}

type Provider interface {
	Authenticator
	Authorizer
}

type Manager interface {
	CreateUser(ctx context.Context, name string) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	IssueToken(ctx context.Context, user string) (string, error)
	RevokeTokens(ctx context.Context, user string) error
	Grant(ctx context.Context, user, collection string, level Level) error
	Ungrant(ctx context.Context, user, collection string) error
}
