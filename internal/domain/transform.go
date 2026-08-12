package domain

type Transform struct {
	Name string
	Version string
	Params map[string]string
}

type Processor interface {
	Process(data []byte, transform Transform) ([]byte, error)
	Name() string
	Version() string
}

func (t Transform) Key() string {
	return t.Name + "@" + t.Version
}

