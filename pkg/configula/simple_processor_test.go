package configula

import (
	"bytes"
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestRecursiveGenerate(t *testing.T) {
	tests := []struct {
		node     *yaml.Node
		expected string
	}{
		{&yaml.Node{Tag: "!!int", Value: "1"}, "YamlNode(1)"},
		{&yaml.Node{Tag: "!!str", Value: "foobar"}, "YamlNode('foobar')"},
		{
			&yaml.Node{
				Tag: "!!map",
				Content: []*yaml.Node{
					&yaml.Node{Tag: "!!str", Value: "foo"},
					&yaml.Node{Tag: "!!str", Value: "bar"},
					&yaml.Node{Tag: "!!str", Value: "baz"},
					&yaml.Node{Tag: "!!int", Value: "2"},
				},
			},
			"{'foo':YamlNode('bar'), 'baz':YamlNode(2), }",
		},
		{
			&yaml.Node{
				Tag: "!!seq",
				Content: []*yaml.Node{
					&yaml.Node{Tag: "!!str", Value: "foo"},
					&yaml.Node{Tag: "!!str", Value: "bar"},
					&yaml.Node{Tag: "!!str", Value: "baz"},
					&yaml.Node{Tag: "!!int", Value: "2"},
				},
			},
			"[YamlNode('foo'),YamlNode('bar'),YamlNode('baz'),YamlNode(2)]",
		},
	}
	for _, test := range tests {
		output := recursiveGenerate("", test.node)
		if output != test.expected {
			t.Errorf("expected %s saw %s", test.expected, output)
		}
	}
}

func TestSimpleProcessor(t *testing.T) {
	tests := []struct {
		data     string
		sections []Section
	}{
		{
			data: "foo = bar: baz",
			sections: []Section{
				{
					Data:      []byte("bar: baz"),
					LineStart: Position{1, 6},
					LineEnd:   Position{1, 14},
					Yaml:      "YamlVariable({'bar':YamlNode('baz'), })",
				},
			},
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
					LineEnd:   Position{13, 4},
					Yaml:      "YamlVariable({'apiVersion':YamlNode('v1'), 'kind':YamlNode('Namespace'), 'metadata':{'name':YamlExpr(lambda: user), }, })",
				},
			},
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
    blaht:
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
	blaht:
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
					LineEnd:   Position{23, 0},
					Yaml:      "YamlVariable({'foo':YamlNode('bar bar'), 'baz':YamlNode('blah'), 'bar':YamlNode(1), 'bob':{'foo':YamlNode('bar'), 'blaht':{'blah':YamlNode('blahber'), }, 'someList':[YamlNode('a'),YamlNode('b'),YamlNode('c')], 'otherList':[YamlNode(1),YamlNode(2),YamlNode(3)], 'anotherList':[YamlNode('a'),YamlNode('b'),YamlNode('c'),YamlExpr(lambda: c + 'd' + fn('bazzer'))], }, })",
				},
				Section{
					Data:      []byte("foo: !~ a + b"),
					LineStart: Position{25, 0},
					LineEnd:   Position{26, 0},
					Yaml:      "YamlVariable({'foo':YamlExpr(lambda: a + b), })",
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
					LineEnd:   Position{33, 0},
					Yaml:      "YamlVariable({'alist':[YamlNode('one'),YamlExpr(lambda: 1 + 4),YamlNode('two'),YamlNode('three')], })",
				},
			},
		},
	}
	for _, test := range tests {
		parser := NewSimpleParser()
		processor := NewSimpleProcessor()
		_, sections, err := parser.GetSections(bytes.NewBufferString(test.data))
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(sections) != len(test.sections) {
			t.Errorf("Unexpected sections:\n%#v\n%#v", sections, test.sections)
		}
		if err := processor.Process(sections); err != nil {
			t.Errorf("%v", err)
		}
		for ix := range test.sections {
			expected := test.sections[ix]
			actual := sections[ix]
			if !reflect.DeepEqual(expected, actual) {
				//t.Errorf("Mismatch [case %d]:\n%#v\n%#v", ix, expected, actual)
				if expected.Yaml != actual.Yaml {
					t.Errorf("Data:\n%s\n%s", expected.Yaml, actual.Yaml)
				}
				if expected.LineStart.Line != actual.LineStart.Line ||
					expected.LineStart.Character != actual.LineStart.Character {
					t.Errorf("Start mismatch: %#v vs %#v", expected.LineStart, actual.LineStart)
				}
				if expected.LineEnd.Line != actual.LineEnd.Line ||
					expected.LineEnd.Character != actual.LineEnd.Character {
					t.Errorf("Start mismatch: %#v vs %#v", expected.LineEnd, actual.LineEnd)
				}
			}
		}
	}
}
