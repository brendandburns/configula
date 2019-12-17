package configula

import (
	"fmt"
	"strings"
	yaml "gopkg.in/yaml.v3"
)

type simpleProcessor struct {}

func NewSimpleProcessor() Processor {
	return &simpleProcessor{}
}

func genKeyValues(indent string, nodes []*yaml.Node) string {
	ix := 0
	result := ""
	for ix < len(nodes) {
		result += fmt.Sprintf("%s'%s':", indent, nodes[ix].Value)
		if (ix + 1 >= len(nodes)) {
			panic(fmt.Sprintf("ERR: %#v\n", nodes[ix]))
		}
		result += recursiveGenerate(indent, nodes[ix+1])
		result += ",\n"
		ix += 2
	}
	return result
}

func recursiveGenerate(indent string, node *yaml.Node) string {
	switch(node.Tag) {
	case "":
		result := ""
		for ix := range node.Content {
			result += recursiveGenerate("", node.Content[ix])
		}
		return result
	case "!!map":
		result := fmt.Sprintf("%s{\n", indent)
		if (len(node.Content)) > 0 {
			result += genKeyValues(indent, node.Content);
		}
		result += indent + "}\n"
		return result;
	case "!!str":
		return fmt.Sprintf("%s'%s'", indent, node.Value)
	case "!!seq":
		items := []string{}
		for ix := range node.Content {
			items = append(items, recursiveGenerate(indent, node.Content[ix]))
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ","))
	case "!~":
		return node.Value
	default:
		return "<custom>"
	}
}

func (s *simpleProcessor) Process(sections []Section) error {
	for ix := range sections {
		node := yaml.Node{}
		if err := yaml.Unmarshal(sections[ix].Data, &node); err != nil {
			return err
		}
		sections[ix].Yaml = recursiveGenerate("", &node)
	}
	return nil
}