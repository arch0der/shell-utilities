// yaml2json - Convert YAML to JSON (simple subset: scalar, list, map)
// Usage: yaml2json [file]
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	var r io.Reader = os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "yaml2json:", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}
	sc := bufio.NewScanner(r)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	val := parseYAML(lines, 0)
	b, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "yaml2json:", err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}

func indent(line string) int {
	n := 0
	for _, c := range line {
		if c == ' ' {
			n++
		} else {
			break
		}
	}
	return n
}

func parseYAML(lines []string, baseIndent int) interface{} {
	if len(lines) == 0 {
		return nil
	}
	// Check first non-empty line for type
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		ind := indent(line)
		if ind < baseIndent {
			break
		}
		if strings.HasPrefix(trimmed, "- ") || trimmed == "-" {
			// List
			return parseList(lines, ind)
		}
		if strings.Contains(trimmed, ": ") || strings.HasSuffix(trimmed, ":") {
			// Map
			return parseMap(lines[i:], ind)
		}
		// Scalar
		return coerce(trimmed)
	}
	return nil
}

func parseList(lines []string, baseInd int) []interface{} {
	var result []interface{}
	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			i++
			continue
		}
		ind := indent(line)
		if ind < baseInd {
			break
		}
		if ind == baseInd && strings.HasPrefix(trimmed, "- ") {
			val := strings.TrimPrefix(trimmed, "- ")
			// Collect sub-lines
			j := i + 1
			for j < len(lines) {
				if strings.TrimSpace(lines[j]) == "" {
					j++
					continue
				}
				if indent(lines[j]) <= baseInd {
					break
				}
				j++
			}
			if indent_check := strings.TrimSpace(val); indent_check != "" {
				result = append(result, coerce(val))
			} else {
				result = append(result, parseYAML(lines[i+1:j], baseInd+2))
			}
			i = j
		} else {
			i++
		}
	}
	return result
}

func parseMap(lines []string, baseInd int) map[string]interface{} {
	result := map[string]interface{}{}
	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			i++
			continue
		}
		ind := indent(line)
		if ind < baseInd {
			break
		}
		if ind == baseInd && strings.Contains(trimmed, ":") {
			colonIdx := strings.Index(trimmed, ":")
			key := trimmed[:colonIdx]
			val := strings.TrimSpace(trimmed[colonIdx+1:])
			// Collect sub-lines
			j := i + 1
			for j < len(lines) {
				if strings.TrimSpace(lines[j]) == "" {
					j++
					continue
				}
				if indent(lines[j]) <= baseInd {
					break
				}
				j++
			}
			if val != "" {
				result[key] = coerce(val)
			} else {
				result[key] = parseYAML(lines[i+1:j], baseInd+2)
			}
			i = j
		} else {
			i++
		}
	}
	return result
}

func coerce(s string) interface{} {
	if s == "null" || s == "~" {
		return nil
	}
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	// Strip quotes
	if (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) ||
		(strings.HasPrefix(s, `'`) && strings.HasSuffix(s, `'`)) {
		return s[1 : len(s)-1]
	}
	return s
}
