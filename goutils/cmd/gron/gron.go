// gron - Make JSON greppable by converting to path=value assignments.
//
// Usage:
//
//	gron [OPTIONS] [FILE...]
//	cat data.json | gron
//	cat data.json | gron | grep email | gron --ungron
//
// Options:
//
//	-u, --ungron  Convert gron output back to JSON
//	-c            Compact JSON output (with --ungron)
//	-s SEP        Assignment separator (default: " = ")
//	--no-sort     Don't sort keys
//
// Examples:
//
//	curl -s api/users | gron | grep '"admin"'
//	gron data.json | grep 'user\[0\]'
//	gron data.json | grep name | gron -u    # filtered JSON
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	ungron  = flag.Bool("u", false, "ungron: convert back to JSON")
	compact = flag.Bool("c", false, "compact output with --ungron")
	noSort  = flag.Bool("no-sort", false, "don't sort keys")
)

func gronify(path string, v interface{}, lines *[]string) {
	switch val := v.(type) {
	case map[string]interface{}:
		*lines = append(*lines, fmt.Sprintf("%s = {};", path))
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		if !*noSort {
			sort.Strings(keys)
		}
		for _, k := range keys {
			subPath := path + "." + k
			// quote keys with special chars
			if strings.ContainsAny(k, ". -[]()") {
				subPath = path + `["` + k + `"]`
			}
			gronify(subPath, val[k], lines)
		}
	case []interface{}:
		*lines = append(*lines, fmt.Sprintf("%s = [];", path))
		for i, item := range val {
			gronify(fmt.Sprintf("%s[%d]", path, i), item, lines)
		}
	default:
		b, _ := json.Marshal(v)
		*lines = append(*lines, fmt.Sprintf("%s = %s;", path, string(b)))
	}
}

func setPath(root map[string]interface{}, path []string, val interface{}) {
	if len(path) == 0 {
		return
	}
	if len(path) == 1 {
		root[path[0]] = val
		return
	}
	key := path[0]
	next, ok := root[key]
	if !ok {
		next = make(map[string]interface{})
		root[key] = next
	}
	if m, ok := next.(map[string]interface{}); ok {
		setPath(m, path[1:], val)
	}
}

func parseGronLine(line string) ([]string, interface{}) {
	// json.key[0].sub = VALUE;
	eq := strings.Index(line, " = ")
	if eq < 0 {
		return nil, nil
	}
	pathStr := line[:eq]
	valStr := strings.TrimSuffix(strings.TrimSpace(line[eq+3:]), ";")

	// parse path
	pathStr = strings.TrimPrefix(pathStr, "json")
	var parts []string
	for pathStr != "" {
		if pathStr[0] == '.' {
			pathStr = pathStr[1:]
			i := strings.IndexAny(pathStr, ".[")
			if i < 0 {
				parts = append(parts, pathStr)
				break
			}
			parts = append(parts, pathStr[:i])
			pathStr = pathStr[i:]
		} else if pathStr[0] == '[' {
			j := strings.Index(pathStr, "]")
			if j < 0 {
				break
			}
			idx := pathStr[1:j]
			idx = strings.Trim(idx, `"`)
			parts = append(parts, idx)
			pathStr = pathStr[j+1:]
		} else {
			break
		}
	}

	var val interface{}
	json.Unmarshal([]byte(valStr), &val)
	return parts, val
}

func processGron(r *os.File) {
	dec := json.NewDecoder(bufio.NewReader(r))
	var data interface{}
	if err := dec.Decode(&data); err != nil {
		fmt.Fprintf(os.Stderr, "gron: %v\n", err)
		os.Exit(1)
	}
	var lines []string
	gronify("json", data, &lines)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for _, l := range lines {
		fmt.Fprintln(w, l)
	}
}

func processUngron(r *os.File) {
	root := make(map[string]interface{})
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || line == "json = {};" || line == "json = [];" {
			continue
		}
		parts, val := parseGronLine(line)
		if parts != nil {
			setPath(root, parts, val)
		}
	}
	var b []byte
	if *compact {
		b, _ = json.Marshal(root)
	} else {
		b, _ = json.MarshalIndent(root, "", "  ")
	}
	fmt.Println(string(b))
}

func main() {
	flag.BoolVar(ungron, "ungron", false, "ungron")
	flag.Parse()
	files := flag.Args()

	process := processGron
	if *ungron {
		process = processUngron
	}

	if len(files) == 0 {
		process(os.Stdin)
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "gron: %v\n", err)
				os.Exit(1)
			}
			process(fh)
			fh.Close()
		}
	}

	_ = strconv.Itoa // avoid unused
}
