// jsonformat - pretty-print or minify JSON
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	minify := false
	indent := "  "
	color := false
	files := []string{}

	for _, arg := range os.Args[1:] {
		switch {
		case arg == "-m" || arg == "--minify": minify = true
		case arg == "-c" || arg == "--color": color = true
		case strings.HasPrefix(arg, "--indent="): indent = arg[len("--indent="):]
		case strings.HasPrefix(arg, "-i"): indent = arg[2:]
		default: files = append(files, arg)
		}
	}

	process := func(r io.Reader, name string) {
		data, err := io.ReadAll(r)
		if err != nil { fmt.Fprintln(os.Stderr, err); return }
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			fmt.Fprintf(os.Stderr, "jsonformat: %s: %v\n", name, err); return
		}
		var out []byte
		if minify {
			out, _ = json.Marshal(v)
		} else {
			out, _ = json.MarshalIndent(v, "", indent)
		}
		if color {
			fmt.Println(colorJSON(string(out)))
		} else {
			os.Stdout.Write(append(out, '\n'))
		}
	}

	if len(files) == 0 { process(os.Stdin, "<stdin>"); return }
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		process(fh, f); fh.Close()
	}
}

func colorJSON(s string) string {
	const (
		reset  = "\033[0m"
		key    = "\033[34m"   // blue keys
		str    = "\033[32m"   // green strings
		num    = "\033[33m"   // yellow numbers
		bool_  = "\033[35m"   // magenta booleans
		null_  = "\033[31m"   // red null
	)
	var b strings.Builder
	inStr := false
	isKey := true
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch == '"' && !inStr:
			inStr = true
			if isKey { b.WriteString(key + `"`) } else { b.WriteString(str + `"`) }
		case ch == '"' && inStr:
			inStr = false; b.WriteByte(ch); b.WriteString(reset); isKey = false
		case inStr: b.WriteByte(ch)
		case ch == ':': b.WriteByte(ch); isKey = false
		case ch == ',' || ch == '{' || ch == '}' || ch == '[' || ch == ']':
			b.WriteByte(ch); if ch == ',' || ch == '{' || ch == '[' { isKey = true }
		case ch == 't' && i+3 < len(s) && s[i:i+4] == "true": b.WriteString(bool_+"true"+reset); i += 3
		case ch == 'f' && i+4 < len(s) && s[i:i+5] == "false": b.WriteString(bool_+"false"+reset); i += 4
		case ch == 'n' && i+3 < len(s) && s[i:i+4] == "null": b.WriteString(null_+"null"+reset); i += 3
		case (ch >= '0' && ch <= '9') || ch == '-':
			b.WriteString(num); b.WriteByte(ch)
			for i+1 < len(s) && (s[i+1] >= '0' && s[i+1] <= '9' || s[i+1] == '.' || s[i+1] == 'e' || s[i+1] == 'E' || s[i+1] == '+' || s[i+1] == '-') { i++; b.WriteByte(s[i]) }
			b.WriteString(reset)
		default: b.WriteByte(ch)
		}
	}
	return b.String()
}
