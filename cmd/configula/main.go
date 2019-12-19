package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/brendandburns/configula/pkg/configula"
	flag "github.com/spf13/pflag"
)

var (
	pythonCommand = flag.String("python", "python3", "The executable to run for Python, overridden by $CONFIGULA_PYTHON")
	dryRun        = flag.Bool("debug", false, "If true, only output the interim program, don't execute")
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: configula [--debug] [--python=/some/path] <path/to/config/file>\n")
		os.Exit(-1)
	}
	file, err := os.Open(flag.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read %s: %s", flag.Args()[0], err.Error())
		os.Exit(1)
	}

	parser := configula.NewSimpleParser()
	lines, sections, err := parser.GetSections(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse %s: %s", flag.Args()[0], err.Error())
		os.Exit(2)
	}
	if err = file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close %s: %s", flag.Args()[0], err.Error())
		os.Exit(3)
	}

	processor := configula.NewSimpleProcessor()
	err = processor.Process(sections)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to process file %s: %s", flag.Args()[0], err.Error())
		os.Exit(4)
	}

	generator := configula.NewPythonGenerator()
	reader, err := generator.Generate(lines, sections)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate: %s", err.Error())
		os.Exit(1)
	}
	if *dryRun {
		if _, err := io.Copy(os.Stdout, reader); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy data %s", err.Error())
		}
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
		fmt.Fprintf(os.Stderr, "Failed to execute Python %s: %s", pythonExec, err.Error())
		os.Exit(1)
	}
}
