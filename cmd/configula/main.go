package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/configula/configula/pkg/configula"
	flag "github.com/spf13/pflag"
)

var (
	pythonCommand = flag.String("python", "python3", "The executable to run for Python, overridden by $CONFIGULA_PYTHON")
	dryRun = flag.Bool("debug", false, "If true, only output the interim program, don't execute")
)

func main() {
	flag.Parse()
	parser := configula.NewSimpleParser()
	file, err := os.Open(flag.Args()[0])
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

	generator := configula.NewPythonGenerator()
	reader, err := generator.Generate(lines, sections)
	if err != nil {
		panic(err)
	}
	if *dryRun {
		io.Copy(os.Stdout, reader)
		return
	}

	pythonExec := *pythonCommand
	if len(os.Getenv("CONFIGULA_PYTHON")) > 0 {
		pythonExec = os.Getenv("CONFIGULA_PYTHON")
	}

	cmd := exec.Command(pythonExec)
	cmd.Stdin = reader
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		panic(err)
	}
}
