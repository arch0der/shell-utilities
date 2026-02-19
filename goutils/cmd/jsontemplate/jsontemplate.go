// jsontemplate - apply a JSON object as template variables to a text template
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: jsontemplate <data.json> [template_file]
  Apply JSON data to a Go text template.
  If no template file, reads template from stdin.
  Template syntax: {{.key}}  {{if .flag}}...{{end}}  {{range .list}}...{{end}}`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 { usage() }
	dataFile := os.Args[1]

	dataBytes, err := os.ReadFile(dataFile)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	var data interface{}
	if err := json.Unmarshal(dataBytes, &data); err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

	var tmplBytes []byte
	if len(os.Args) > 2 {
		tmplBytes, err = os.ReadFile(os.Args[2])
	} else {
		tmplBytes, err = io.ReadAll(os.Stdin)
	}
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

	funcMap := template.FuncMap{
		"upper": strings.ToUpper, "lower": strings.ToLower,
		"trim": strings.TrimSpace, "replace": strings.ReplaceAll,
		"join": func(sep string, s []interface{}) string {
			parts := make([]string, len(s))
			for i, v := range s { parts[i] = fmt.Sprintf("%v", v) }
			return strings.Join(parts, sep)
		},
		"default": func(def, val interface{}) interface{} {
			if val == nil || val == "" { return def }; return val
		},
	}

	t, err := template.New("t").Funcs(funcMap).Parse(string(tmplBytes))
	if err != nil { fmt.Fprintln(os.Stderr, "parse error:", err); os.Exit(1) }
	if err := t.Execute(os.Stdout, data); err != nil {
		fmt.Fprintln(os.Stderr, "render error:", err); os.Exit(1)
	}
}
