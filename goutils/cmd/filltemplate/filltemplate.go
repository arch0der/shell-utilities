// filltemplate - fill {{PLACEHOLDER}} style templates from env, flags, or JSON
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var placeholderRe = regexp.MustCompile(`\{\{([^}]+)\}\}`)

func fill(tmpl string, vars map[string]string, missing string) (string, []string) {
	var unresolved []string
	result := placeholderRe.ReplaceAllStringFunc(tmpl, func(m string) string {
		key := strings.TrimSpace(m[2:len(m)-2])
		if val, ok := vars[key]; ok { return val }
		if missing == "empty" { return "" }
		if missing == "keep" { return m }
		unresolved = append(unresolved, key)
		return m
	})
	return result, unresolved
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: filltemplate [options] [template_file]
  -v KEY=VALUE   set a variable (repeatable)
  -j <json>      JSON object of key→value pairs (or @file)
  -e             include environment variables
  --missing      how to handle missing vars: keep|empty|error (default: keep)
  If no template file given, reads from stdin.`)
	os.Exit(1)
}

func main() {
	vars := map[string]string{}
	useEnv := false
	missing := "keep"
	var tmplFile string
	var jsonData string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-v":
			i++; parts := strings.SplitN(args[i], "=", 2)
			if len(parts) == 2 { vars[parts[0]] = parts[1] }
		case "-j": i++; jsonData = args[i]
		case "-e": useEnv = true
		case "--missing": i++; missing = args[i]
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			tmplFile = args[i]
		}
	}

	if useEnv {
		for _, e := range os.Environ() {
			p := strings.SplitN(e, "=", 2)
			if len(p) == 2 { vars[p[0]] = p[1] }
		}
	}

	if jsonData != "" {
		raw := jsonData
		if strings.HasPrefix(raw, "@") {
			b, err := os.ReadFile(raw[1:]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
			raw = string(b)
		}
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &m); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		for k, v := range m { vars[k] = fmt.Sprintf("%v", v) }
	}

	var tmplBytes []byte
	var err error
	if tmplFile != "" { tmplBytes, err = os.ReadFile(tmplFile) } else { tmplBytes, err = io.ReadAll(os.Stdin) }
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

	result, unresolved := fill(string(tmplBytes), vars, missing)
	fmt.Print(result)
	if len(unresolved) > 0 && missing == "error" {
		fmt.Fprintf(os.Stderr, "filltemplate: unresolved placeholders: %s\n", strings.Join(unresolved, ", "))
		os.Exit(1)
	}
}
