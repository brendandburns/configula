package configula

import (
	"io"
)

type Position struct {
	Line int
	Character int
}

type Section struct {
	Data []byte
	LineStart Position
	LineEnd Position
	Yaml string
	HasYaml bool
}

type Parser interface {
	GetSections(io.Reader) ([]Section, error)
}