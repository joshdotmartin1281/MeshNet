package domain

type Source string

const (
	SourceCLI    Source = "cli"
	SourceHTTP   Source = "http"
	SourceMANUAL Source = "manual"
)
