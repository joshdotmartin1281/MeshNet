package text

import (
	"fmt"
	"strings"

	"MeshNet/internal/domain"
)

type Uppercase struct{}

func NewUppercase() Uppercase {
	return Uppercase{}
}

func (Uppercase) Name() string {
	return "uppercase"
}

func (Uppercase) Version() string {
	return "1"
}

func (Uppercase) Process(data []byte, transform domain.Transform) ([]byte, error) {
	if transform.Key() != "uppercase@1" {
		return nil, fmt.Errorf(
			"unsupported transform: %s",
			transform.Key(),
		)
	}

	return []byte(strings.ToUpper(string(data))), nil
}
