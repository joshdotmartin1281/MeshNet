package text

import (
	"fmt"
	"strings"

	"MeshNet/transform"
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

func (Uppercase) Process(data []byte, transform transform.Transform) ([]byte, error) {
	if transform.Key() != "uppercase@1" {
		return nil, fmt.Errorf(
			"unsupported transform: %s",
			transform.Key(),
		)
	}

	return []byte(strings.ToUpper(string(data))), nil
}
