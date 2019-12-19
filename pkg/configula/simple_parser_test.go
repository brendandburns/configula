package configula

import (
	"testing"
)

func TestFindSpace(t *testing.T) {
	tests := []struct {
		value string
		start int
		expected int
		name string
	} {
		{ "abc1234:asbd", 10, -1, "no space" },
		{ "abc 123 : acd", -1, -1, "negative start"},
		{ "abc 123 : acd", 5, 3, "simple" },
		{ "abc\t123:acd", 10, 3, "tabs" },
	}
	for _, test := range tests {
		if ix := findSpace(test.value, test.start); ix != test.expected {
			t.Errorf("[%s]: Expected %d, saw %d", test.name, test.expected, ix)
		}	
	}
}

func TestExtractYaml(t *testing.T) {
	tests := []struct {
		value string
		start int
		end int
		yaml string
		name string
	} {
		{ "foo = baz: bar", 5, 14, "baz: bar", "simple" }, 
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