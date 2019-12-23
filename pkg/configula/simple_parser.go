package configula

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type simpleParser struct{}

func extractYaml(str string) (int, int, string) {
	regex := regexp.MustCompile(`^(?:.*[\s\(])*([a-z0-9]+:\s*(?:(?:[a-z0-9]+)|(?:!.*)))(?:[\s\)].*)*$`)
	index := regex.FindSubmatchIndex([]byte(str))
	if index != nil {
		return index[2], index[3], str[index[2]:index[3]]
	}
	return -1, -1, str
}

func isPythonLine(str string) bool {
	if strings.Index(str, ":") == -1 {
		return true
	}
	return strings.HasPrefix(str, "def") ||
		strings.HasPrefix(str, "if") ||
		strings.HasPrefix(str, "elif") ||
		strings.HasPrefix(str, "else") ||
		strings.HasPrefix(str, "while") ||
		strings.HasPrefix(str, "for") ||
		strings.HasPrefix(str, "class") ||
		strings.HasPrefix(str, "try") ||
		strings.HasPrefix(str, "except") ||
		strings.Contains(str, "lambda")
}

func isYamlStart(line string) int {
	if ix := strings.Index(line, ":"); ix != -1 {
		start, _, _ := extractYaml(line)
		if start == 2 {
			return start
		}
	}
	return -1
}

func isYamlEnd(line string) int {
	trim := strings.TrimSpace(line)
	if len(trim) == 0 {
		return 0
	}
	return -1
}

func trimBrackets(str string) string {
	start := strings.Index(str, "<")
	end := strings.LastIndex(str, ">")
	if start == -1 && end == -1 {
		return str
	}
	if start == -1 {
		return str[0:end]
	}
	if end == -1 {
		return str[start+1:]
	}
	return str[start+1:end]
}

func (*simpleParser) GetSections(reader io.Reader) ([]string, []Section, error) {
	scanner := bufio.NewScanner(reader)
	result := []Section{}
	lineNum := 0
	parenCount := 0
	blockStart := Position{-1, -1}
	yamlStart := Position{-1, -1}
	blockBuffer := ""
	lines := []string{}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		lines = append(lines, line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if ix := strings.Index(line, "<"); ix != -1 {
			if parenCount == 0 {
				blockStart.Line = lineNum
				blockStart.Character = ix
			}
			parenCount++
			blockBuffer += line + "\n"
			continue
		}
		if blockStart.Line != -1 {
			blockBuffer += line + "\n"
			if ix := strings.Index(line, ">"); ix != -1 {
				parenCount--
				if parenCount == 0 {
					if blockStart.Line != lineNum {
						end := Position{lineNum, ix}
						fmt.Printf("Adding here!\n")
						trimmed := trimBrackets(blockBuffer)
						fmt.Printf("%s\n%s\n", blockBuffer, trimmed)
						result = append(result, Section{[]byte(trimmed), blockStart, end, ""})
					}
					blockStart = Position{-1, -1}
					blockBuffer = ""
				}
			}
			continue
		}
		if isPythonLine(strings.TrimSpace(line)) {
			if yamlStart.Line != -1 {
				end := Position{lineNum, 0}
				result = append(result, Section{[]byte(blockBuffer), yamlStart, end, "" })
				yamlStart = Position{-1, -1}
				blockBuffer = ""
			}
			continue
		}
		if yamlStart.Line != -1 {
			blockBuffer += line + "\n"
			if ix := isYamlEnd(line); ix != -1 {
				end := Position{lineNum, 0}
				result = append(result, Section{[]byte(blockBuffer), yamlStart, end, "" })
				yamlStart = Position{-1, -1}
				blockBuffer = ""
			}
			continue
		}
		if ix := isYamlStart(line); ix != -1 {
			yamlStart.Line = lineNum
			yamlStart.Character = ix
			blockBuffer = line + "\n"
			continue
		}
		start, end, data := extractYaml(line)
		result = append(result,
			Section{[]byte(data), Position{lineNum, start}, Position{lineNum, end}, ""})
	}
	return lines, result, nil
}

// NewSimpleParser creates a parser that is simple
func NewSimpleParser() Parser {
	return &simpleParser{}
}
