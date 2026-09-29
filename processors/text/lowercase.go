package text

import (
	"fmt"
	"strings"

	"MeshNet/transform"
)

type Lowercase struct{}

func NewLowercase() Lowercase {
	return Lowercase{}
}

func (Lowercase) Name() string {
	return "lowercase"
}

func (Lowercase) Version() string {
	return "1"
}

func (Lowercase) Process(data []byte, transform transform.Transform) ([]byte, error) {
	if transform.Key() != "lowercase@1" {
		return nil, fmt.Errorf(
			"unsupported transform: %s",
			transform.Key(),
		)
	}

	return []byte(strings.ToLower(string(data))), nil
}
