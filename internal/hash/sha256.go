package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

type SHA256 struct{}

func NewSHA256() SHA256 {
	return SHA256{}
}

func (SHA256) Hash(data []byte) string {
	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:])
}

func (SHA256) Name() string {
	return "sha256"
}

func (SHA256) Version() string {
	return "1"
}
