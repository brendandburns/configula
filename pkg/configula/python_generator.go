package configula

import (
	"bytes"
	"fmt"
	"io"
)

type pythonGenerator struct {}

func NewPythonGenerator() Generator {
	return &pythonGenerator{}
}

func (p *pythonGenerator) Generate(lines []string, sections []Section) (io.Reader, error) {
	lineNum := 1
	section := 0
	buff := &bytes.Buffer{}
	if _, err := fmt.Fprintf(buff, "from runtime.configula import (render, YamlExpr, YamlNode, YamlVariable)\n"); err != nil {
		return nil, err
	}
	
	for lineNum < len(lines) + 1 && section < len(sections) {
		currentSection := sections[section]
		line := lines[lineNum - 1]
		if lineNum < currentSection.LineStart.Line {
			if _, err := fmt.Fprintf(buff, "%s\n", line); err != nil {
				return nil, err
			}
		}
		if lineNum == currentSection.LineStart.Line {
			if currentSection.LineStart.Character != 0 {
				if _, err := fmt.Fprintf(buff, line[0:currentSection.LineStart.Character]); err != nil {
					return nil, err
				}
			}
			if _, err := fmt.Fprintf(buff, currentSection.Yaml); err != nil {
				return nil, err
			}
			lineNum = currentSection.LineEnd.Line
			line := lines[lineNum - 1]
			if currentSection.LineEnd.Character + 1 < len(line) {
				if _, err := fmt.Fprintf(buff, "%s\n", line[currentSection.LineEnd.Character + 1:]); err != nil {
					return nil, err
				}
			}

			section++
		}
		lineNum++
	}
	for ; lineNum < len(lines) + 1; lineNum++ {
		line := lines[lineNum - 1]
		if _, err := fmt.Fprintf(buff, "%s\n", line); err != nil {
			return nil, err
		}
	}
	return buff, nil
}