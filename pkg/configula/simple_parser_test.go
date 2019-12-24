package configula

import (
	"bytes"
	"reflect"
	"testing"
)

func TestIsPythonLine(t * testing.T) {
	tests := []struct{
		line string
		isPython bool
	} {
		{ "foo = 1 + 2", true},
		{ "def fn():", true},
		{ "bar: baz", false},
		{ "lambda x: x * 2", true},
		{ "if foo:", true},
		{ "else: ", true},
		{ "elif foo:", true},
	}

	for _, test := range tests {
		if isPythonLine(test.line) {
			if !test.isPython {
				t.Errorf("Expected %s to not be python but it was", test.line)
			}
		} else {
			if test.isPython {
				t.Errorf("Expected %s to be Python but it's not", test.line)
			}
		}
	}
}

func TestExtractYaml(t *testing.T) {
	tests := []struct {
		value string
		start int
		end   int
		yaml  string
		name  string
	}{
		{"foo = baz: bar", 6, 14, "baz: bar", "simple"},
		{"return baz: bar", 7, 15, "baz: bar", "return"},
		{"render(baz: bar)", 7, 15, "baz: bar", "function"},
		{"foo = baz: !~ 1 + 2", 6, 19, "baz: !~ 1 + 2", "yaml tag"},
		// TODO {"foo: [1, 2, 3]", 0, 5, "foo: [1, 2, 3]", "yaml sequence"},
	}
	for _, test := range tests {
		start, end, yaml := extractYaml(test.value)
		if start != test.start {
			t.Errorf("[%s]: Expected %d, saw %d", test.name, test.start, start)
		}
		if end != test.end {
			t.Errorf("[%s]: Expected %d, saw %d", test.name, test.end, end)
		}
		if yaml != test.yaml {
			t.Errorf("[%s]: Expected '%s', saw '%s'", test.name, test.yaml, yaml)
		}
	}
}

func TestSimpleParser(t *testing.T) {
	tests := []struct {
		data     string
		sections []Section
		name string
	}{
		{
			data: "foo = bar: baz",
			sections: []Section{
				{
					Data: []byte("bar: baz"),
					LineStart: Position{1, 6},
					LineEnd:   Position{1, 14},
				},
			},
			name: "simple",
		},
		{
			data: `
foo =
  bar: baz
  blah: bar

foo.render()
`,
			sections: []Section{
				{
					Data: []byte(
`  bar: baz
  blah: bar
`),
					LineStart: Position{3, 2},
					LineEnd:   Position{5, 0},
				},
			},
			name: "no brackets",
		},
		{
			data: `
users = ['jim', 'sally', 'sue']

ns =
  apiVersion: v1
  kind: Namespace
  metadata:
    name: !~ user

for user in users:
  ns.render()
`,
			sections: []Section{
				{
					Data: []byte(
`  apiVersion: v1
  kind: Namespace
  metadata:
    name: !~ user
`),
					LineStart: Position{5, 2},
					LineEnd:   Position{9, 0},
				},
			},
			name: "bigger no brackets",
		},
		{
			data: `
# Simple example of creating 3 Kubernetes namespaces

# Our users in need of namespaces
users = ['jim', 'sally', 'sue']

# The namespaces objects from YAML
namespaces = map(lambda user: <
        apiVersion: v1
        kind: Namespace
        metadata:
          name: !~ user
    >, users)

# Output
render(namespaces)
`,
			sections: []Section{
				Section{
					Data: []byte(`
        apiVersion: v1
        kind: Namespace
        metadata:
          name: !~ user
    `),
					LineStart: Position{8, 30},
					LineEnd: Position{13, 4},
				},
			},
			name: "brackets",
		},
		{
			data: `foobar = 'blah %d' % 1
c = 'a'
			
def fn(val):
  return val + '!!'
			
baz = <
  foo: "bar bar"
  baz: blah
  bar: 1
  bob:
	foo: bar
	blah:
	  blah: blahber
	someList:
	- a
	- b
	- c
	otherList: [1, 2, 3]
	anotherList: [
	  a, b, c, !~ c + 'd' + fn('bazzer')
	]
>
			
foo: !~ a + b
			
<
  alist:
  - one
  - !~ 1 + 4
  - two
  - three
>
			
baz.render()`,
			sections: []Section{
				Section{
					Data: []byte(`
  foo: "bar bar"
  baz: blah
  bar: 1
  bob:
	foo: bar
	blah:
	  blah: blahber
	someList:
	- a
	- b
	- c
	otherList: [1, 2, 3]
	anotherList: [
	  a, b, c, !~ c + 'd' + fn('bazzer')
	]
`),
					LineStart: Position{7, 6},
					LineEnd: Position{23, 0},
				},
				Section{
					Data: []byte("foo: !~ a + b\n"),
					LineStart: Position{25, 0},
					LineEnd: Position{26, 0},
				},
				Section{
					Data: []byte(`
  alist:
  - one
  - !~ 1 + 4
  - two
  - three
`),
					LineStart: Position{27, 0},
					LineEnd: Position{33, 0},
				},
			},
			name: "complex",
		},
	}
	for _, test := range tests {
		parser := NewSimpleParser()
		_, sections, err := parser.GetSections(bytes.NewBufferString(test.data))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(sections) != len(test.sections) {
			t.Errorf("[Test %s] Expected %d saw %d sections", test.name, len(test.sections), len(sections))
			t.Errorf("[Test %s] Unexpected sections:\n%#v\n%#v", test.name, sections, test.sections)
			for ix := range sections {
				t.Errorf("%s\n", string(sections[ix].Data))
			}
			t.FailNow()
		}
		for ix := range test.sections {
			expected := test.sections[ix]
			actual := sections[ix]
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("[Test %s ] Unexpected section (%d):\n%#v\n%#v", test.name, ix, string(expected.Data), string(actual.Data))
				t.Errorf("[Test %s ] Unexpected section (%d):\n%#v\n%#v", test.name, ix, expected, actual)
			}
		}
	}
}
