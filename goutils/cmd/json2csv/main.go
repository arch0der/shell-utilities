// json2csv - Convert JSON array of objects to CSV
// Usage: json2csv [-d delim] [-no-header] [file]
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

var (
	delim    = flag.String("d", ",", "Field delimiter")
	noHeader = flag.Bool("no-header", false, "Omit header row")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: json2csv [-d delim] [-no-header] [file]")
		flag.PrintDefaults()
	}
	flag.Parse()

	var r io.Reader = os.Stdin
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "json2csv:", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}

	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "json2csv:", err)
		os.Exit(1)
	}

	var input interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintln(os.Stderr, "json2csv: invalid JSON:", err)
		os.Exit(1)
	}

	w := csv.NewWriter(os.Stdout)
	if len(*delim) > 0 {
		w.Comma = []rune(*delim)[0]
	}
	defer w.Flush()

	switch v := input.(type) {
	case []interface{}:
		if len(v) == 0 {
			return
		}
		// If array of objects
		if obj, ok := v[0].(map[string]interface{}); ok {
			// Collect all keys
			keySet := map[string]bool{}
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					for k := range m {
						keySet[k] = true
					}
				}
			}
			keys := make([]string, 0, len(keySet))
			for k := range keySet {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			_ = obj

			if !*noHeader {
				w.Write(keys)
			}
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					row := make([]string, len(keys))
					for i, k := range keys {
						row[i] = fmt.Sprintf("%v", m[k])
						if m[k] == nil {
							row[i] = ""
						}
					}
					w.Write(row)
				}
			}
		} else {
			// Array of arrays or primitives
			for _, item := range v {
				switch row := item.(type) {
				case []interface{}:
					strs := make([]string, len(row))
					for i, cell := range row {
						strs[i] = fmt.Sprintf("%v", cell)
					}
					w.Write(strs)
				default:
					w.Write([]string{fmt.Sprintf("%v", item)})
				}
			}
		}
	case map[string]interface{}:
		if !*noHeader {
			w.Write([]string{"key", "value"})
		}
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			w.Write([]string{k, fmt.Sprintf("%v", v[k])})
		}
	default:
		fmt.Fprintln(os.Stderr, "json2csv: expected array or object")
		os.Exit(1)
	}
}
