package domain

type Hasher interface {
	Hash(data []byte) string
	Name() string
	Version() string
}

