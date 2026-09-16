package domain

import (
	"fmt"
	"strings"
)

type Transform struct {
	Name    string
	Version string
	Params  map[string]string
}

type Processor interface {
	Process(data []byte, transform Transform) ([]byte, error)
	Name() string
	Version() string
}

func (t Transform) Key() string {
	return t.Name + "@" + t.Version
}

func ParseTransform(value string) (Transform, error) {
	value = strings.TrimSpace(value)

	parts := strings.SplitN(value, "@", 2)

	if len(parts) != 2 {
		return Transform{}, fmt.Errorf(
			"invalid transform %q: expected name@version",
			value,
		)
	}

	name := strings.TrimSpace(parts[0])
	version := strings.TrimSpace(parts[1])

	if name == "" || version == "" {
		return Transform{}, fmt.Errorf(
			"invalid transform %q: expected name@version",
			value,
		)
	}

	return Transform{
		Name:    name,
		Version: version,
		Params:  map[string]string{},
	}, nil
}
