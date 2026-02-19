// json2yaml - Convert JSON to YAML.
//
// Usage:
//
//	json2yaml [FILE...]
//	cat data.json | json2yaml
//
// Options:
//
//	-i N      Indent spaces (default: 2)
//
// Examples:
//
//	cat config.json | json2yaml
//	json2yaml package.json
//	echo '{"name":"alice","age":30}' | json2yaml
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

var indent = flag.Int("i", 2, "indent spaces")

func toYAML(v interface{}, depth int) string {
	pad := strings.Repeat(" ", depth * *indent)
	childPad := strings.Repeat(" ", (depth+1) * *indent)
	_ = childPad

	switch val := v.(type) {
	case nil:
		return "null"
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		if val == math.Trunc(val) && !math.IsInf(val, 0) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case string:
		// Quote strings that could be misinterpreted
		needsQuote := false
		if val == "" || val == "true" || val == "false" || val == "null" ||
			val == "yes" || val == "no" || val == "on" || val == "off" {
			needsQuote = true
		}
		if strings.ContainsAny(val, ":#{}[]|>&*!,") || strings.HasPrefix(val, "- ") {
			needsQuote = true
		}
		if needsQuote {
			escaped := strings.ReplaceAll(val, `"`, `\"`)
			return fmt.Sprintf(`"%s"`, escaped)
		}
		if strings.Contains(val, "\n") {
			lines := strings.Split(val, "\n")
			var sb strings.Builder
			sb.WriteString("|\n")
			for _, l := range lines {
				sb.WriteString(childPad + l + "\n")
			}
			return strings.TrimRight(sb.String(), "\n")
		}
		return val
	case map[string]interface{}:
		if len(val) == 0 {
			return "{}"
		}
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var sb strings.Builder
		for _, k := range keys {
			child := val[k]
			childStr := toYAML(child, depth+1)
			switch child.(type) {
			case map[string]interface{}, []interface{}:
				if child != nil {
					sb.WriteString(fmt.Sprintf("%s%s:\n%s", pad, k, childStr))
				} else {
					sb.WriteString(fmt.Sprintf("%s%s: null\n", pad, k))
				}
			default:
				sb.WriteString(fmt.Sprintf("%s%s: %s\n", pad, k, childStr))
			}
		}
		return sb.String()
	case []interface{}:
		if len(val) == 0 {
			return "[]\n"
		}
		var sb strings.Builder
		for _, item := range val {
			childStr := toYAML(item, depth+1)
			switch item.(type) {
			case map[string]interface{}:
				lines := strings.SplitAfter(childStr, "\n")
				if len(lines) > 0 && lines[0] != "" {
					sb.WriteString(fmt.Sprintf("%s- %s", pad, strings.TrimPrefix(lines[0], pad)))
					for _, l := range lines[1:] {
						if l != "" {
							sb.WriteString(l)
						}
					}
				}
			default:
				sb.WriteString(fmt.Sprintf("%s- %s\n", pad, childStr))
			}
		}
		return sb.String()
	}
	return fmt.Sprintf("%v", v)
}

func process(r *os.File) {
	dec := json.NewDecoder(bufio.NewReader(r))
	var data interface{}
	if err := dec.Decode(&data); err != nil {
		fmt.Fprintf(os.Stderr, "json2yaml: %v\n", err)
		os.Exit(1)
	}
	out := toYAML(data, 0)
	fmt.Print("---\n" + out)
}

func main() {
	flag.Parse()
	files := flag.Args()
	if len(files) == 0 {
		process(os.Stdin)
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "json2yaml: %v\n", err)
			os.Exit(1)
		}
		process(fh)
		fh.Close()
	}
}
