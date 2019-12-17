package main

import (
	//"fmt"
	"bufio"
	"os"

	"github.com/configula/configula/pkg/configula"
)

func main() {
	parser := configula.NewSimpleParser()
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	sections, err := parser.GetSections(file)
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

	file, err = os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
}