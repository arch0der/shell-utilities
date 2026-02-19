// jsonpath - Query JSON using dot-notation paths.
//
// Usage:
//
//	jsonpath [OPTIONS] EXPR [FILE...]
//	cat data.json | jsonpath .users[0].name
//
// Options:
//
//	-r        Raw output (no quotes for strings)
//	-c        Compact output (no pretty-print)
//	-e        Exit with error if path not found
//	-a        Iterate array at root; apply expr to each element
//
// Expressions:
//
//	.key           object field
//	.key.sub       nested field
//	.[N]           array index (0-based)
//	.key[N].sub    combined
//	.              root (pretty-print the whole thing)
//
// Examples:
//
//	echo '{"name":"alice","age":30}' | jsonpath .name        # "alice"
//	echo '{"name":"alice","age":30}' | jsonpath -r .name     # alice
//	cat users.json | jsonpath .users[0].email
//	cat users.json | jsonpath -a .name    # print .name from each element
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	raw     = flag.Bool("r", false, "raw string output")
	compact = flag.Bool("c", false, "compact output")
	errExit = flag.Bool("e", false, "error if not found")
	array   = flag.Bool("a", false, "iterate root array")
)

func lookup(v interface{}, path []string) (interface{}, bool) {
	if len(path) == 0 {
		return v, true
	}
	part := path[0]
	rest := path[1:]

	// array index: [N]
	if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
		idx, err := strconv.Atoi(part[1 : len(part)-1])
		if err != nil {
			return nil, false
		}
		arr, ok := v.([]interface{})
		if !ok || idx < 0 || idx >= len(arr) {
			return nil, false
		}
		return lookup(arr[idx], rest)
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, false
	}
	val, ok := m[part]
	if !ok {
		return nil, false
	}
	return lookup(val, rest)
}

func parsePath(expr string) []string {
	// expr: .foo.bar[0].baz  → ["foo","bar","[0]","baz"]
	expr = strings.TrimPrefix(expr, ".")
	if expr == "" {
		return nil
	}
	// split on . and [
	var parts []string
	for _, seg := range strings.Split(expr, ".") {
		// handle embedded [N]
		for {
			i := strings.Index(seg, "[")
			if i < 0 {
				if seg != "" {
					parts = append(parts, seg)
				}
				break
			}
			if i > 0 {
				parts = append(parts, seg[:i])
			}
			j := strings.Index(seg[i:], "]")
			if j < 0 {
				parts = append(parts, seg)
				break
			}
			parts = append(parts, seg[i:i+j+1])
			seg = seg[i+j+1:]
		}
	}
	return parts
}

func printVal(v interface{}) {
	if v == nil {
		fmt.Println("null")
		return
	}
	switch val := v.(type) {
	case string:
		if *raw {
			fmt.Println(val)
		} else {
			b, _ := json.Marshal(val)
			fmt.Println(string(b))
		}
	default:
		var b []byte
		if *compact {
			b, _ = json.Marshal(v)
		} else {
			b, _ = json.MarshalIndent(v, "", "  ")
		}
		fmt.Println(string(b))
	}
}

func query(data interface{}, path []string) {
	if *array {
		arr, ok := data.([]interface{})
		if !ok {
			fmt.Fprintln(os.Stderr, "jsonpath: -a requires root array")
			os.Exit(1)
		}
		for _, item := range arr {
			val, found := lookup(item, path)
			if !found {
				if *errExit {
					fmt.Fprintln(os.Stderr, "jsonpath: path not found")
					os.Exit(1)
				}
				continue
			}
			printVal(val)
		}
		return
	}
	val, found := lookup(data, path)
	if !found {
		if *errExit {
			fmt.Fprintln(os.Stderr, "jsonpath: path not found")
			os.Exit(1)
		}
		return
	}
	printVal(val)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: jsonpath EXPR [FILE...]")
		os.Exit(1)
	}

	expr := args[0]
	path := parsePath(expr)
	files := args[1:]

	decode := func(r *os.File) {
		dec := json.NewDecoder(bufio.NewReader(r))
		var data interface{}
		if err := dec.Decode(&data); err != nil {
			fmt.Fprintf(os.Stderr, "jsonpath: %v\n", err)
			os.Exit(1)
		}
		query(data, path)
	}

	if len(files) == 0 {
		decode(os.Stdin)
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "jsonpath: %v\n", err)
				os.Exit(1)
			}
			decode(fh)
			fh.Close()
		}
	}
}
