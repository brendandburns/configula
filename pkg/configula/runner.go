package configula

import "io"

type Runner struct {
	Parser Parser
	Proc Processor
	Gen Generator
	Exec Executor
}

func (r *Runner) Run(input io.Reader, output io.Writer, dryRun bool) error {
	lines, sections, err := r.Parser.GetSections(input)
	if err != nil {
		return err
	}

	err = r.Proc.Process(sections)
	if err != nil {
		return err
	}

	reader, err := r.Gen.Generate(lines, sections)
	if err != nil {
		return err
	}
	if dryRun {
		_, err := io.Copy(output, reader)
		return err
	}

	return r.Exec.Execute(output, reader)
}