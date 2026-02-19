// template - render Go text templates with env vars or JSON data
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
	fmt.Fprintln(os.Stderr, `usage: template [options] <template_file>
  -d <json>   JSON data (string or @file)
  -e          use environment variables as data
  -s <str>    template string instead of file
  Reads template from file; outputs rendered result.`)
	os.Exit(1)
}

func main() {
	var dataJSON string
	useEnv := false
	var tmplStr string
	var tmplFile string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d": i++; dataJSON = args[i]
		case "-e": useEnv = true
		case "-s": i++; tmplStr = args[i]
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			if tmplFile == "" { tmplFile = args[i] }
		}
	}

	// Load template
	var tmplContent string
	if tmplStr != "" {
		tmplContent = tmplStr
	} else if tmplFile != "" {
		b, err := os.ReadFile(tmplFile)
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		tmplContent = string(b)
	} else {
		b, _ := io.ReadAll(os.Stdin)
		tmplContent = string(b)
	}

	// Build data map
	data := map[string]interface{}{}
	if useEnv {
		for _, e := range os.Environ() {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 { data[parts[0]] = parts[1] }
		}
	}
	if dataJSON != "" {
		var raw string
		if strings.HasPrefix(dataJSON, "@") {
			b, err := os.ReadFile(dataJSON[1:])
			if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
			raw = string(b)
		} else { raw = dataJSON }
		var extra map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &extra); err != nil {
			fmt.Fprintln(os.Stderr, "template: bad JSON:", err); os.Exit(1)
		}
		for k, v := range extra { data[k] = v }
	}

	funcMap := template.FuncMap{
		"upper": strings.ToUpper, "lower": strings.ToLower,
		"title": strings.Title, "trim": strings.TrimSpace,
		"replace": strings.ReplaceAll, "split": strings.Split,
		"join": strings.Join,
		"env": os.Getenv,
	}

	t, err := template.New("t").Funcs(funcMap).Parse(tmplContent)
	if err != nil { fmt.Fprintln(os.Stderr, "template parse error:", err); os.Exit(1) }
	if err := t.Execute(os.Stdout, data); err != nil {
		fmt.Fprintln(os.Stderr, "template exec error:", err); os.Exit(1)
	}
}
