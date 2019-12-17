package main

import (
	"fmt"
	"os"

	"github.com/configula/configula/pkg/configula"
)

func main() {
	parser := configula.NewSimpleParser()
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	lines, sections, err := parser.GetSections(file)
	if err != nil {
		panic(err)
	}
	if err = file.Close(); err != nil {
		panic(err)
	}

	processor := configula.NewSimpleProcessor()
	err = processor.Process(sections)
	if err != nil {
		panic(err)
	}

	lineNum := 1
	section := 0
	fmt.Printf("from runtime.configula import (YamlExpr, YamlNode, YamlVariable)\n")
	for lineNum < len(lines) + 1 && section < len(sections) {
		currentSection := sections[section]
		line := lines[lineNum - 1]
		if lineNum < currentSection.LineStart.Line {
			fmt.Printf("%s\n", line)
		}
		if lineNum == currentSection.LineStart.Line {
			fmt.Printf("%s%s",
				line[0:currentSection.LineStart.Character - 1],
				currentSection.Yaml)
			lineNum = currentSection.LineEnd.Line
			line := lines[lineNum - 1]
			fmt.Printf("%s\n", line[currentSection.LineEnd.Character + 1:])

			section++
		}
		lineNum++
	}
	for ; lineNum < len(lines) + 1; lineNum++ {
		line := lines[lineNum - 1]
		fmt.Printf("%s\n", line)
	}
}