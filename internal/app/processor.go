package app

import (
	"fmt"

	"MeshNet/internal/domain"
)

type Processor struct {
	processors map[string]domain.Processor
}

func NewProcessor(processors ...domain.Processor) *Processor {
	registry := make(map[string]domain.Processor)

	for _, processor := range processors {
		key := processor.Name() + "@" + processor.Version()

		registry[key] = processor
	}

	return &Processor{
		processors: registry,
	}
}

func (p *Processor) Process(data []byte, transforms []domain.Transform) ([]byte, error) {
	result := data

	for _, transform := range transforms {
		processor, ok := p.processors[transform.Key()]
		if !ok {
			return nil, fmt.Errorf(
				"processor not found: %s",
				transform.Key(),
			)
		}

		var err error

		result, err = processor.Process(result, transform)
		if err != nil {
			return nil, fmt.Errorf(
				"process %s: %w",
				transform.Key(),
				err,
			)
		}
	}

	return result, nil
}
