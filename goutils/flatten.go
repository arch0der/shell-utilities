// flatten - Flatten nested JSON to dot-notation key=value pairs.
//
// Usage:
//
//	flatten [OPTIONS] [FILE...]
//	cat nested.json | flatten
//
// Options:
//
//	-s SEP    Key separator (default: .)
//	-j        Output as flat JSON object (default: key=value lines)
//	-p PREFIX Add prefix to all keys
//	-a        Expand arrays (default: include index in key)
//
// Examples:
//
//	echo '{"a":{"b":1}}'     | flatten         # a.b=1
//	echo '{"a":{"b":1}}'     | flatten -j       # {"a.b":1}
//	echo '{"x":[1,2,3]}'     | flatten          # x.0=1  x.1=2  x.2=3
//	cat config.json          | flatten -s __     # a__b=1
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
	sep    = flag.String("s", ".", "separator")
	asJSON = flag.Bool("j", false, "output as JSON")
	prefix = flag.String("p", "", "key prefix")
)

func flattenVal(prefix string, v interface{}, out map[string]interface{}) {
	switch val := v.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			newKey := k
			if prefix != "" {
				newKey = prefix + *sep + k
			}
			flattenVal(newKey, val[k], out)
		}
	case []interface{}:
		for i, item := range val {
			newKey := fmt.Sprintf("%s%s%d", prefix, *sep, i)
			if prefix == "" {
				newKey = fmt.Sprintf("%d", i)
			}
			flattenVal(newKey, item, out)
		}
	default:
		out[prefix] = v
	}
}

func valStr(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func process(r *os.File) {
	dec := json.NewDecoder(bufio.NewReader(r))
	var data interface{}
	if err := dec.Decode(&data); err != nil {
		fmt.Fprintf(os.Stderr, "flatten: %v\n", err)
		os.Exit(1)
	}

	out := make(map[string]interface{})
	flattenVal(*prefix, data, out)

	if *asJSON {
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return
	}

	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, valStr(out[k]))
	}
}

func main() {
	flag.Parse()
	files := flag.Args()

	_ = strings.Contains // avoid unused import

	if len(files) == 0 {
		process(os.Stdin)
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "flatten: %v\n", err)
				os.Exit(1)
			}
			process(fh)
			fh.Close()
		}
	}
}
