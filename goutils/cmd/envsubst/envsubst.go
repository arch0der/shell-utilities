// envsubst - Substitute environment variables in text.
//
// Usage:
//
//	envsubst [OPTIONS] [FILE...]
//	cat template.txt | envsubst
//
// Options:
//
//	-v VARS   Only substitute specified vars (comma-separated)
//	-d        Fail if any referenced variable is unset
//	-n        Don't substitute, just list referenced variables
//	-e KEY=VAL  Set extra variable (repeatable, before substitution)
//	-p        Also expand ${var:-default} and ${var:+alt} syntax
//
// Substitutions supported:
//   $VAR         simple variable
//   ${VAR}       braced variable
//   ${VAR:-DEF}  default if unset/empty  (with -p)
//   ${VAR:+ALT}  alternate if set       (with -p)
//
// Examples:
//
//	cat k8s.yaml | envsubst
//	NAME=world envsubst <<< "Hello, $NAME!"
//	envsubst -v DATABASE_URL,PORT config.template
//	envsubst -e "ENV=prod" -e "VERSION=1.2" template.txt
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type extraVars []string

func (e *extraVars) String() string  { return strings.Join(*e, ",") }
func (e *extraVars) Set(v string) error { *e = append(*e, v); return nil }

var extras extraVars

var (
	onlyVars = flag.String("v", "", "only these vars")
	strict   = flag.Bool("d", false, "fail on unset")
	listOnly = flag.Bool("n", false, "list variables")
	extended = flag.Bool("p", false, "extended syntax")
)

var (
	simpleVar   = regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
	bracedVar   = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)
	defaultVar  = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*):-([^}]*)\}`)
	altVar      = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*):([+])([^}]*)\}`)
)

func main() {
	flag.Var(&extras, "e", "extra variable KEY=VAL")
	flag.Parse()

	// Build env map
	env := make(map[string]string)
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	for _, e := range extras {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	var allowedVars map[string]bool
	if *onlyVars != "" {
		allowedVars = make(map[string]bool)
		for _, v := range strings.Split(*onlyVars, ",") {
			allowedVars[strings.TrimSpace(v)] = true
		}
	}

	allowed := func(name string) bool {
		if allowedVars == nil {
			return true
		}
		return allowedVars[name]
	}

	foundVars := make(map[string]bool)

	subst := func(s string) string {
		// Extended: ${VAR:-default} and ${VAR:+alt}
		if *extended {
			s = defaultVar.ReplaceAllStringFunc(s, func(m string) string {
				parts := defaultVar.FindStringSubmatch(m)
				name, def := parts[1], parts[2]
				foundVars[name] = true
				if !allowed(name) {
					return m
				}
				v, ok := env[name]
				if !ok || v == "" {
					return def
				}
				return v
			})
			s = altVar.ReplaceAllStringFunc(s, func(m string) string {
				parts := altVar.FindStringSubmatch(m)
				name, alt := parts[1], parts[3]
				foundVars[name] = true
				if !allowed(name) {
					return m
				}
				if _, ok := env[name]; ok {
					return alt
				}
				return ""
			})
		}

		// ${VAR}
		s = bracedVar.ReplaceAllStringFunc(s, func(m string) string {
			name := bracedVar.FindStringSubmatch(m)[1]
			foundVars[name] = true
			if !allowed(name) {
				return m
			}
			v, ok := env[name]
			if !ok {
				if *strict {
					fmt.Fprintf(os.Stderr, "envsubst: unset variable: %s\n", name)
					os.Exit(1)
				}
				return ""
			}
			return v
		})

		// $VAR
		s = simpleVar.ReplaceAllStringFunc(s, func(m string) string {
			name := simpleVar.FindStringSubmatch(m)[1]
			foundVars[name] = true
			if !allowed(name) {
				return m
			}
			v, ok := env[name]
			if !ok {
				if *strict {
					fmt.Fprintf(os.Stderr, "envsubst: unset variable: %s\n", name)
					os.Exit(1)
				}
				return ""
			}
			return v
		})

		return s
	}

	files := flag.Args()
	var readers []*os.File
	if len(files) == 0 {
		readers = []*os.File{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "envsubst: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}

	if *listOnly {
		// scan without substituting
		for _, r := range readers {
			sc := bufio.NewScanner(r)
			for sc.Scan() {
				subst(sc.Text())
			}
		}
		for v := range foundVars {
			fmt.Println(v)
		}
		return
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for _, r := range readers {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1024*1024), 1024*1024)
		for sc.Scan() {
			fmt.Fprintln(w, subst(sc.Text()))
		}
	}
}
