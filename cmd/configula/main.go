package main

import (
	"fmt"
	"os/exec"
	"io"
	"strings"
	"os"

	"github.com/brendandburns/configula/pkg/configula"
	flag "github.com/spf13/pflag"
)

var (
	pythonCommand = flag.String("python", "python3", "The executable to run for Python, overridden by $CONFIGULA_PYTHON")
	dryRun        = flag.Bool("debug", false, "If true, only output the interim program, don't execute")
	file = flag.StringP("filename", "f", "", "The file name to process")
)

func pluginMain() {
	if len(*file) == 0 || len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: kubectl configula create|apply|delete -f <some-file>\n")
		return
	}
	cmd := exec.Command("kubectl", flag.Args()[0], "-f", "-")
	output, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		process(*file, output)
		if err := output.Close(); err != nil {
			panic(err)
		}
	}()
	cmd.Run()
}

func main() {
	flag.Parse()
	if strings.HasSuffix(os.Args[0], "kubectl-configula") {
		pluginMain()
		return
	}
	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: configula [--debug] [--python=/some/path] <path/to/config/file>\n")
		os.Exit(-1)
	}
	process(flag.Args()[0], os.Stdout)
}

func process(filename string, output io.Writer) {
	file, err := os.Open(filename)
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

	if err := rn.Run(file, output, *dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute: %s\n", err)
		os.Exit(-1)
	}
}
