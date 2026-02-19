// globmatch - test glob patterns against filenames or filter stdin lines
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: globmatch <pattern> [strings...]  |  globmatch <pattern> < lines
  Test whether strings match a glob pattern.
  -v    invert match (print non-matching)
  -q    quiet: exit 0 if any match, 1 if none
  -f    match against actual filesystem (expand glob)
  Patterns: * ? [abc] [a-z] {a,b,c} (brace expansion)`)
	os.Exit(1)
}

func bracketExpand(pattern string) []string {
	start := strings.Index(pattern, "{")
	if start < 0 { return []string{pattern} }
	end := strings.Index(pattern[start:], "}") + start
	if end <= start { return []string{pattern} }
	prefix := pattern[:start]
	suffix := pattern[end+1:]
	alts := strings.Split(pattern[start+1:end], ",")
	var result []string
	for _, alt := range alts {
		for _, sub := range bracketExpand(prefix+alt+suffix) { result = append(result, sub) }
	}
	return result
}

func matchAny(patterns []string, s string) bool {
	for _, p := range patterns {
		matched, err := filepath.Match(p, s)
		if err == nil && matched { return true }
		// also try matching just the filename part
		matched, err = filepath.Match(p, filepath.Base(s))
		if err == nil && matched { return true }
	}
	return false
}

func main() {
	invert := false
	quiet := false
	fsMode := false
	args := os.Args[1:]
	filtered := args[:0]
	for _, a := range args {
		switch a {
		case "-v": invert = true
		case "-q": quiet = true
		case "-f": fsMode = true
		default: filtered = append(filtered, a)
		}
	}
	args = filtered
	if len(args) < 1 { usage() }

	pattern := args[0]
	patterns := bracketExpand(pattern)

	if fsMode {
		for _, p := range patterns {
			matches, err := filepath.Glob(p)
			if err != nil { fmt.Fprintln(os.Stderr, err); continue }
			for _, m := range matches { fmt.Println(m) }
		}
		return
	}

	matched := false
	process := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" { return }
		m := matchAny(patterns, s)
		if invert { m = !m }
		if m {
			matched = true
			if !quiet { fmt.Println(s) }
		}
	}

	if len(args) > 1 {
		for _, s := range args[1:] { process(s) }
	} else {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() { process(sc.Text()) }
	}

	if quiet { if !matched { os.Exit(1) } }
}
