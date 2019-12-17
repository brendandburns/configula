package configula

import (
	"io"
)

// Position represents a (Line, Char) position in a stream
type Position struct {
	Line int
	Character int
}

// Section represents a YAML section
type Section struct {
	Data []byte
	LineStart Position
	LineEnd Position
	Yaml string
}

// Parser is an interface for all parsers
type Parser interface {
	// GetSections returns the set of lines in the file and the YAML sections
	GetSections(io.Reader) ([]string, []Section, error)
}