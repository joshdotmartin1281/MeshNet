package encrypt

import (
	"fmt"

	"MeshNet/internal/domain"
)

type Decrypt struct {
	key []byte
}

func NewDecrypt(key []byte) Decrypt {
	return Decrypt{key: key}
}

func (Decrypt) Name() string {
	return "decrypt"
}

func (Decrypt) Version() string {
	return "1"
}

func (d Decrypt) Process(data []byte, transform domain.Transform) ([]byte, error) {
	if transform.Key() != "decrypt@1" {
		return nil, fmt.Errorf(
			"unsupported transform: %s",
			transform.Key(),
		)
	}

	aead, err := NewGCM(d.key)
	if err != nil {
		return nil, err
	}

	nonceSize := aead.NonceSize()

	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt data: %w", err)
	}

	return plaintext, nil
}
