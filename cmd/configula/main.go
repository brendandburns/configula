package main

import (
	"fmt"
	"os"

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
	defer file.Close()

	pythonExec := *pythonCommand
	if len(os.Getenv("CONFIGULA_PYTHON")) > 0 {
		pythonExec = os.Getenv("CONFIGULA_PYTHON")
	}

	rn := &configula.Runner{
		configula.NewSimpleParser(),
		configula.NewSimpleProcessor(),
		configula.NewPythonGenerator(),
		configula.NewPythonExecutor(pythonExec),
	}

	if err := rn.Run(file, os.Stdout, *dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute: %s", err)
		os.Exit(-1)
	}
}
