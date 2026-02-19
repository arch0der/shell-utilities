// fromjson - Pretty-print, validate, and inspect JSON.
//
// Usage:
//
//	fromjson [OPTIONS] [FILE...]
//	echo '{"a":1}' | fromjson
//
// Options:
//
//	-c        Compact (minify) output
//	-s        Sort keys
//	-t        Print type summary
//	-k        Print all keys (flattened dot-notation)
//	-v        Validate only (exit 0=valid, 1=invalid)
//	-d N      Max depth to display (truncate deep objects)
//	-C        Colorize output (requires terminal)
//
// Examples:
//
//	cat messy.json | fromjson              # pretty print
//	cat data.json | fromjson -c            # minify
//	cat data.json | fromjson -s            # sorted keys
//	echo '{"a":1' | fromjson -v            # validate
//	cat big.json | fromjson -k            # list all keys
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	compact    = flag.Bool("c", false, "compact")
	sortKeys   = flag.Bool("s", false, "sort keys")
	typeSummary = flag.Bool("t", false, "type summary")
	listKeys   = flag.Bool("k", false, "list keys")
	validateOnly = flag.Bool("v", false, "validate only")
	maxDepth   = flag.Int("d", 0, "max depth")
	colorize   = flag.Bool("C", false, "colorize")
)

const (
	colorReset  = "\033[0m"
	colorKey    = "\033[34m"    // blue
	colorString = "\033[32m"    // green
	colorNum    = "\033[33m"    // yellow
	colorNull   = "\033[90m"    // gray
	colorBool   = "\033[35m"    // magenta
)

func sortedMarshal(v interface{}, indent string, depth int) string {
	switch val := v.(type) {
	case nil:
		if *colorize {
			return colorNull + "null" + colorReset
		}
		return "null"
	case bool:
		s := "false"
		if val {
			s = "true"
		}
		if *colorize {
			return colorBool + s + colorReset
		}
		return s
	case float64:
		s := fmt.Sprintf("%g", val)
		if *colorize {
			return colorNum + s + colorReset
		}
		return s
	case string:
		b, _ := json.Marshal(val)
		if *colorize {
			return colorString + string(b) + colorReset
		}
		return string(b)
	case []interface{}:
		if len(val) == 0 {
			return "[]"
		}
		if *maxDepth > 0 && depth >= *maxDepth {
			return fmt.Sprintf("[... %d items]", len(val))
		}
		child := indent + "  "
		var parts []string
		for _, item := range val {
			parts = append(parts, child+sortedMarshal(item, child, depth+1))
		}
		return "[\n" + strings.Join(parts, ",\n") + "\n" + indent + "]"
	case map[string]interface{}:
		if len(val) == 0 {
			return "{}"
		}
		if *maxDepth > 0 && depth >= *maxDepth {
			return fmt.Sprintf("{... %d keys}", len(val))
		}
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		if *sortKeys {
			sort.Strings(keys)
		}
		child := indent + "  "
		var parts []string
		for _, k := range keys {
			kb, _ := json.Marshal(k)
			keyStr := string(kb)
			if *colorize {
				keyStr = colorKey + keyStr + colorReset
			}
			parts = append(parts, child+keyStr+": "+sortedMarshal(val[k], child, depth+1))
		}
		return "{\n" + strings.Join(parts, ",\n") + "\n" + indent + "}"
	}
	return "null"
}

func gatherKeys(v interface{}, prefix string, out map[string]string) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, child := range val {
			p := k
			if prefix != "" {
				p = prefix + "." + k
			}
			gatherKeys(child, p, out)
		}
	case []interface{}:
		for i, child := range val {
			p := fmt.Sprintf("%s[%d]", prefix, i)
			gatherKeys(child, p, out)
		}
	default:
		t := fmt.Sprintf("%T", val)
		if val == nil {
			t = "null"
		}
		out[prefix] = t
	}
}

func process(r *os.File) bool {
	dec := json.NewDecoder(bufio.NewReader(r))
	var data interface{}
	if err := dec.Decode(&data); err != nil {
		fmt.Fprintf(os.Stderr, "fromjson: invalid JSON: %v\n", err)
		return false
	}

	if *validateOnly {
		return true
	}

	if *listKeys {
		keys := make(map[string]string)
		gatherKeys(data, "", keys)
		sorted := make([]string, 0, len(keys))
		for k := range keys {
			sorted = append(sorted, k)
		}
		sort.Strings(sorted)
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()
		for _, k := range sorted {
			fmt.Fprintf(w, "%-40s %s\n", k, keys[k])
		}
		return true
	}

	if *typeSummary {
		typeCounts := make(map[string]int)
		var walk func(v interface{})
		walk = func(v interface{}) {
			switch val := v.(type) {
			case map[string]interface{}:
				typeCounts["object"]++
				for _, child := range val {
					walk(child)
				}
			case []interface{}:
				typeCounts["array"]++
				for _, child := range val {
					walk(child)
				}
			case string:
				typeCounts["string"]++
			case float64:
				typeCounts["number"]++
			case bool:
				typeCounts["boolean"]++
			case nil:
				typeCounts["null"]++
			}
		}
		walk(data)
		for k, n := range typeCounts {
			fmt.Printf("%-10s %d\n", k+":", n)
		}
		return true
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if *compact {
		b, _ := json.Marshal(data)
		fmt.Fprintln(w, string(b))
	} else if *sortKeys || *colorize || *maxDepth > 0 {
		fmt.Fprintln(w, sortedMarshal(data, "", 0))
	} else {
		b, _ := json.MarshalIndent(data, "", "  ")
		fmt.Fprintln(w, string(b))
	}
	return true
}

func main() {
	flag.Parse()
	files := flag.Args()

	ok := true
	if len(files) == 0 {
		if !process(os.Stdin) {
			ok = false
		}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fromjson: %v\n", err)
				ok = false
				continue
			}
			if !process(fh) {
				ok = false
			}
			fh.Close()
		}
	}

	if *validateOnly {
		if ok {
			fmt.Println("valid")
			os.Exit(0)
		}
		os.Exit(1)
	}
	if !ok {
		os.Exit(1)
	}
}
