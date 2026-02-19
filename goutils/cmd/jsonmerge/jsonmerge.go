// jsonmerge - deep merge multiple JSON files/objects
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func deepMerge(base, override interface{}) interface{} {
	bMap, bOk := base.(map[string]interface{})
	oMap, oOk := override.(map[string]interface{})
	if bOk && oOk {
		result := map[string]interface{}{}
		for k, v := range bMap { result[k] = v }
		for k, v := range oMap {
			if bv, exists := result[k]; exists {
				result[k] = deepMerge(bv, v)
			} else { result[k] = v }
		}
		return result
	}
	// Arrays: override wins by default
	return override
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: jsonmerge <file1.json> <file2.json> [file3.json...]
  Deep-merge JSON objects. Later files override earlier values.
  -a    arrays: append instead of replace
  -p    pretty print output (default: compact)`)
	os.Exit(1)
}

func main() {
	pretty := false
	var files []string
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-p", "--pretty": pretty = true
		default:
			if strings.HasPrefix(arg, "-") { usage() }
			files = append(files, arg)
		}
	}
	if len(files) < 1 { usage() }

	var result interface{}
	for _, f := range files {
		var r io.Reader
		if f == "-" { r = os.Stdin } else {
			fh, err := os.Open(f); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
			defer fh.Close(); r = fh
		}
		data, _ := io.ReadAll(r)
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil { fmt.Fprintf(os.Stderr, "jsonmerge: %s: %v\n", f, err); os.Exit(1) }
		if result == nil { result = v } else { result = deepMerge(result, v) }
	}

	var out []byte
	if pretty { out, _ = json.MarshalIndent(result, "", "  ") } else { out, _ = json.Marshal(result) }
	fmt.Println(string(out))
}
