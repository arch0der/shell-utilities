// linefilter - filter stdin lines using multiple conditions
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type filter struct {
	re      *regexp.Regexp
	invert  bool
	minLen  int
	maxLen  int
	numOnly bool
	blank   bool // include blank
	uniq    bool
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: linefilter [options]
  -m <pat>     include lines matching regex
  -v <pat>     exclude lines matching regex
  -min <n>     minimum line length
  -max <n>     maximum line length
  -n           only lines that are valid numbers
  -B           drop blank/whitespace-only lines
  -u           unique lines only (remove duplicates)
  -i           case-insensitive matching
  Multiple -m and -v flags are ANDed together.`)
	os.Exit(1)
}

func main() {
	var includes, excludes []*regexp.Regexp
	minLen, maxLen := -1, -1
	numOnly := false
	dropBlank := false
	uniq := false
	caseI := false

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-m":
			i++; prefix := ""; if caseI { prefix = "(?i)" }
			re, err := regexp.Compile(prefix + args[i])
			if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
			includes = append(includes, re)
		case "-v":
			i++; prefix := ""; if caseI { prefix = "(?i)" }
			re, err := regexp.Compile(prefix + args[i])
			if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
			excludes = append(excludes, re)
		case "-min": i++; minLen, _ = strconv.Atoi(args[i])
		case "-max": i++; maxLen, _ = strconv.Atoi(args[i])
		case "-n": numOnly = true
		case "-B": dropBlank = true
		case "-u": uniq = true
		case "-i": caseI = true
		default: usage()
		}
	}

	seen := map[string]bool{}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		if dropBlank && strings.TrimSpace(line) == "" { continue }
		if minLen >= 0 && len(line) < minLen { continue }
		if maxLen >= 0 && len(line) > maxLen { continue }
		if numOnly { _, err := strconv.ParseFloat(strings.TrimSpace(line), 64); if err != nil { continue } }
		pass := true
		for _, re := range includes { if !re.MatchString(line) { pass = false; break } }
		if !pass { continue }
		for _, re := range excludes { if re.MatchString(line) { pass = false; break } }
		if !pass { continue }
		if uniq { if seen[line] { continue }; seen[line] = true }
		fmt.Println(line)
	}
}
