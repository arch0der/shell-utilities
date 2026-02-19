// flatten2 - flatten nested JSON arrays/objects to flat key=value pairs
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func flatten(v interface{}, prefix string, sep string, out map[string]interface{}) {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, child := range t {
			newKey := k; if prefix != "" { newKey = prefix + sep + k }
			flatten(child, newKey, sep, out)
		}
	case []interface{}:
		for i, child := range t {
			newKey := fmt.Sprintf("%s%s%d", prefix, sep, i)
			if prefix == "" { newKey = fmt.Sprintf("%d", i) }
			flatten(child, newKey, sep, out)
		}
	default:
		out[prefix] = v
	}
}

func main() {
	sep := "."
	format := "kv"
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-s": i++; sep = args[i]
		case "-f": i++; format = args[i]
		}
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

	out := map[string]interface{}{}
	flatten(v, "", sep, out)

	keys := make([]string, 0, len(out))
	for k := range out { keys = append(keys, k) }
	sort.Strings(keys)

	switch format {
	case "json":
		b, _ := json.MarshalIndent(out, "", "  "); fmt.Println(string(b))
	case "env":
		for _, k := range keys {
			envKey := strings.ToUpper(strings.NewReplacer(".", "_", "-", "_", "[", "_", "]", "").Replace(k))
			fmt.Printf("export %s=%q\n", envKey, fmt.Sprintf("%v", out[k]))
		}
	default: // kv
		for _, k := range keys { fmt.Printf("%s=%v\n", k, out[k]) }
	}
}
