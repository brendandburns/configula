package configula

import (
	"io"
	"os/exec"
)

type pythonExecutor struct {
	path string
}

func NewPythonExecutor(execPath string) Executor {
	return &pythonExecutor{execPath}
}

func (p *pythonExecutor) Execute(output io.Writer, input io.Reader) error {
	cmd := exec.Command(p.path)
	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = output

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}