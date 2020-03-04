package configula

import (
	"fmt"
	"strings"
	"bytes"
	"strconv"
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
		result += ", "
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
		result := fmt.Sprintf("%s{", indent)
		if (len(node.Content)) > 0 {
			result += genKeyValues(indent, node.Content);
		}
		result += indent + "}"
		return result;
	case "!!str":
		return fmt.Sprintf("%sYamlNode(%s)", indent, toPythonString(node.Value))
	case "!!seq":
		items := []string{}
		for ix := range node.Content {
			items = append(items, recursiveGenerate(indent, node.Content[ix]))
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ","))
	case "!~":
		return "YamlExpr(lambda: " + node.Value + ")"
	case "!!int":
		return fmt.Sprintf("%sYamlNode(%s)", indent, node.Value)
	case "!!bool":
		boolVal, _ := strconv.ParseBool(node.Value)
		if (boolVal) {
			return fmt.Sprintf("%sYamlNode(True)", indent)
		} else {
			return fmt.Sprintf("%sYamlNode(False)", indent)
		}
	case "!!null":
		return fmt.Sprintf("%sYamlNode(None)", indent)
	default:
		return fmt.Sprintf("%s", node.Tag)
	}
}

func toPythonString(input string) string {
	var buffer bytes.Buffer
	buffer.WriteByte('\'')

	for _, char := range input {
		switch(char) {
		case '\\':
			buffer.WriteByte('\\')
			buffer.WriteByte('\\')
		case '\'':
			buffer.WriteByte('\\')
			buffer.WriteByte('\'')
		case '\r':
			buffer.WriteByte('\\')
			buffer.WriteByte('r')
		case '\n':
			buffer.WriteByte('\\')
			buffer.WriteByte('n')
		default:
			buffer.WriteRune(char)
		}
	}

	buffer.WriteByte('\'')
	return buffer.String()
}

func (s *simpleProcessor) Process(sections []Section) error {
	for ix := range sections {
		node := yaml.Node{}
		if err := yaml.Unmarshal(sections[ix].Data, &node); err != nil {
			fmt.Printf("%s\n", string(sections[ix].Data))
			return err
		}
		sections[ix].Yaml = "YamlVariable(" + recursiveGenerate("", &node) + ")"
	}
	return nil
}