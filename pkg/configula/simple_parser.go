package configula

import (
	"bufio"
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

func (*simpleParser) GetSections(reader io.Reader) ([]string, []Section, error) {
	scanner := bufio.NewScanner(reader)
	result := []Section{}
	lineNum := 0
	parenCount := 0
	blockStart := Position{-1, -1}
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
		}
		if blockStart.Line != -1 {
			blockBuffer += line + "\n"
			if ix := strings.Index(line, ">"); ix != -1 {
				parenCount--
				if parenCount == 0 {
					if blockStart.Line != lineNum {
						end := Position{lineNum, ix}
						result = append(result, Section{[]byte(blockBuffer), blockStart, end, ""})
					}
					blockStart = Position{-1, -1}
					blockBuffer = ""
				}
			}
			continue
		}
		if isPythonLine(strings.TrimSpace(line)) {
			continue
		}
		pos := Position{lineNum, -1}
		result = append(result, Section{[]byte(line), pos, pos, ""})
	}
	for i := range result {
		str := string(result[i].Data)
		if parenIx := strings.Index(str, "<"); parenIx == -1 {
			start, end, data := extractYaml(str)
			result[i].LineStart.Character = start
			result[i].LineEnd.Character = end
			result[i].Data = []byte(data)
		} else {
			result[i].Data = []byte(str[parenIx+1 : strings.LastIndex(str, ">")])
		}
	}
	return lines, result, nil
}

// NewSimpleParser creates a parser that is simple
func NewSimpleParser() Parser {
	return &simpleParser{}
}
