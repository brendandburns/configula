package configula

import (
	"io"
	"bufio"
	"strings"
)

type simpleParser struct {}

func findSpace(str string, ix int) int {
	for ix >= 0 {
		if str[ix] == ' ' || str[ix] == '\t' {
			return ix
		}
		ix--
		continue
	}
	return -1
}

func extractYaml(str string) (int, int, string) {
	ix := strings.Index(str, ":")
	start := findSpace(str, ix)
	if start != -1 {
		return start + 1, len(str), str[start:]
	}
	return 0, len(str), str
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
		if strings.Index(line, ":") == -1 {
			continue
		}
		strip := strings.TrimSpace(line)
		if strings.HasPrefix(strip, "def") ||
		   strings.HasPrefix(strip, "if") ||
		   strings.HasPrefix(strip, "elif") ||
		   strings.HasPrefix(strip, "else") ||
		   strings.HasPrefix(strip, "while") ||
		   strings.HasPrefix(strip, "for") ||
		   strings.HasPrefix(strip, "class") ||
		   strings.HasPrefix(strip, "try") ||
		   strings.HasPrefix(strip, "except") ||
		   strings.Contains(strip, "lambda") ||
		   strings.HasPrefix(strip, "#") {
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
			result[i].Data = []byte(str[parenIx + 1:strings.LastIndex(str, ">")])
		}
	}
	return lines, result, nil
}

func NewSimpleParser() Parser {
	return &simpleParser{}
}