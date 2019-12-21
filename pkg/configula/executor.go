package configula

import (
	"io"
)

// Executor is a interface for things that execute programs
type Executor interface {
	Execute(output io.Writer, input io.Reader) error
}
