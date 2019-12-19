package configula

import "io"

type Generator interface {
	Generate([]string, []Section) (io.Reader, error)
}