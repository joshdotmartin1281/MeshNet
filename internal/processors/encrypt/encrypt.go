package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"MeshNet/internal/domain"
)

type Encrypt struct {
	key []byte
}

func NewEncrypt(key []byte) Encrypt {
	return Encrypt{key: key}
}

func (Encrypt) Name() string {
	return "encrypt"
}

func (Encrypt) Version() string {
	return "1"
}

func (e Encrypt) Process(data []byte, transform domain.Transform) ([]byte, error) {
	if transform.Key() != "encrypt@1" {
		return nil, fmt.Errorf(
			"unsupported transform: %s",
			transform.Key(),
		)
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := aead.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}
