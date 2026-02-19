// jsonkeys - list all keys/paths in a JSON object (dot-notation)
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func walk(v interface{}, prefix string, paths *[]string) {
	switch t := v.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(t))
		for k := range t { keys = append(keys, k) }
		sort.Strings(keys)
		for _, k := range keys {
			full := k
			if prefix != "" { full = prefix + "." + k }
			*paths = append(*paths, full)
			walk(t[k], full, paths)
		}
	case []interface{}:
		for i, item := range t {
			full := fmt.Sprintf("%s[%d]", prefix, i)
			*paths = append(*paths, full)
			walk(item, full, paths)
		}
	}
}

func main() {
	leafOnly := false
	typed := false
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-l", "--leaf": leafOnly = true
		case "-t", "--type": typed = true
		}
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

	var paths []string
	walk(v, "", &paths)

	typeOf := func(v interface{}) string {
		switch v.(type) {
		case map[string]interface{}: return "object"
		case []interface{}: return "array"
		case float64: return "number"
		case string: return "string"
		case bool: return "bool"
		case nil: return "null"
		}; return "unknown"
	}

	// Build path→value map
	vals := map[string]interface{}{}
	var buildVals func(v interface{}, prefix string)
	buildVals = func(v interface{}, prefix string) {
		vals[prefix] = v
		switch t := v.(type) {
		case map[string]interface{}:
			for k, child := range t {
				full := k; if prefix != "" { full = prefix + "." + k }
				buildVals(child, full)
			}
		case []interface{}:
			for i, item := range t { buildVals(item, fmt.Sprintf("%s[%d]", prefix, i)) }
		}
	}
	buildVals(v, "")

	for _, p := range paths {
		if p == "" { continue }
		val := vals[p]
		isLeaf := true
		switch val.(type) {
		case map[string]interface{}, []interface{}: isLeaf = false
		}
		if leafOnly && !isLeaf { continue }
		line := p
		if typed { line += "  (" + typeOf(val) + ")" }
		if leafOnly && typed {
			s, _ := json.Marshal(val)
			if len(s) > 60 { s = append(s[:57], []byte("...")...) }
			line = fmt.Sprintf("%-50s  =  %s", strings.TrimSpace(p), string(s))
		}
		fmt.Println(line)
	}
}
