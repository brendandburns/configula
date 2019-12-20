package configula

import (
	"bytes"
	"io"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		lines []string
		sections []Section
		output string
		name string
	} {
		{
			lines: []string{
				"foo = bar: baz",
				"foo.render()",
			},
			sections: []Section{
				Section{ 
					Yaml: "YamlModule(bar: baz)",
					LineStart: Position{1, 6},
					LineEnd: Position{1, 13},
				},
			},
			output:
`from runtime.configula import (render, YamlExpr, YamlNode, YamlVariable)
foo = YamlModule(bar: baz)
foo.render()
`,
			name: "simple",
		},
		{
			lines: []string{
				"foo = <",
				"  bar: baz",
				"  blah: flub",
				">",
				"foo.render()",
			},
			sections: []Section{
				Section{ 
					Yaml: "YamlNode(YamlNode(bar: baz),YamlNode(blah: flub))",
					LineStart: Position{1, 6},
					LineEnd: Position{4, 0},
				},
			},
			output:
`from runtime.configula import (render, YamlExpr, YamlNode, YamlVariable)
foo = YamlNode(YamlNode(bar: baz),YamlNode(blah: flub))
foo.render()
`,
			name: "multi-line",
		},
	}
	generator := NewPythonGenerator()
	for _, test := range tests {
		reader, err := generator.Generate(test.lines, test.sections)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		buff := &bytes.Buffer{}
		io.Copy(buff, reader)
		output := buff.String()
		if buff.String() != test.output {
			t.Errorf("[%s] Expected: %s\nSaw: %s\n", test.name, test.output, output)
		}
	}
}