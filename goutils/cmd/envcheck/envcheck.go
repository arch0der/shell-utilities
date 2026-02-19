// envcheck - verify required environment variables are set and meet constraints
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type varSpec struct {
	name     string
	required bool
	pattern  string
	minLen   int
	desc     string
}

func parseSpecFile(path string) ([]varSpec, error) {
	f, err := os.Open(path); if err != nil { return nil, err }
	defer f.Close()
	var specs []varSpec
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") { continue }
		parts := strings.Fields(line)
		spec := varSpec{name: parts[0], required: true}
		for i := 1; i < len(parts); i++ {
			switch {
			case parts[i] == "optional": spec.required = false
			case strings.HasPrefix(parts[i], "pattern="): spec.pattern = parts[i][8:]
			case strings.HasPrefix(parts[i], "minlen="): spec.minLen, _ = strconv.Atoi(parts[i][7:])
			case strings.HasPrefix(parts[i], "desc="): spec.desc = parts[i][5:]
			}
		}
		specs = append(specs, spec)
	}
	return specs, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: envcheck [options] VAR [VAR...]
  -f <file>    read variable specs from file (one per line)
  -p <pattern> regex pattern each var must match
  -r           all vars must be non-empty (default)
  -q           quiet: only exit code
  Spec file format: VAR_NAME [optional] [pattern=regex] [minlen=N] [desc=text]`)
	os.Exit(1)
}

func main() {
	specFile := ""
	pattern := ""
	quiet := false
	var varNames []string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f": i++; specFile = args[i]
		case "-p": i++; pattern = args[i]
		case "-q": quiet = true
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			varNames = append(varNames, args[i])
		}
	}

	var specs []varSpec
	if specFile != "" {
		var err error; specs, err = parseSpecFile(specFile)
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	}
	for _, n := range varNames { specs = append(specs, varSpec{name: n, required: true, pattern: pattern}) }
	if len(specs) == 0 {
		// No args: show all env vars
		for _, e := range os.Environ() { fmt.Println(e) }; return
	}

	allOk := true
	for _, spec := range specs {
		val, set := os.LookupEnv(spec.name)
		ok := true
		var issues []string
		if spec.required && !set { ok = false; issues = append(issues, "not set") } else if spec.required && val == "" { ok = false; issues = append(issues, "empty") }
		if set && val != "" {
			if spec.pattern != "" {
				if matched, _ := regexp.MatchString(spec.pattern, val); !matched {
					ok = false; issues = append(issues, fmt.Sprintf("doesn't match pattern %q", spec.pattern))
				}
			}
			if spec.minLen > 0 && len(val) < spec.minLen {
				ok = false; issues = append(issues, fmt.Sprintf("too short (min %d)", spec.minLen))
			}
		}
		if !quiet {
			status := "✓"; if !ok { status = "✗" }
			desc := ""; if spec.desc != "" { desc = " (" + spec.desc + ")" }
			if ok {
				preview := val; if len(preview) > 30 { preview = preview[:27] + "..." }
				fmt.Printf("%s %-30s = %q%s\n", status, spec.name, preview, desc)
			} else {
				fmt.Printf("%s %-30s [%s]%s\n", status, spec.name, strings.Join(issues, "; "), desc)
			}
		}
		if !ok { allOk = false }
	}
	if !allOk { os.Exit(1) }
}
