// kwsearch - keyword search with context lines and match highlighting
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: kwsearch [options] <pattern> [file...]
  -C <n>    context lines (before and after), default 2
  -B <n>    lines before
  -A <n>    lines after
  -i        case-insensitive
  -n        show line numbers
  -c        count matches only`)
	os.Exit(1)
}

func main() {
	before, after := 2, 2
	caseI, showNums, countOnly := false, true, false
	var pattern string
	var files []string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-C": i++; before, _ = strconv.Atoi(args[i]); after = before
		case "-B": i++; before, _ = strconv.Atoi(args[i])
		case "-A": i++; after, _ = strconv.Atoi(args[i])
		case "-i": caseI = true
		case "-n": showNums = true
		case "-c": countOnly = true
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			if pattern == "" { pattern = args[i] } else { files = append(files, args[i]) }
		}
	}
	if pattern == "" { usage() }

	reStr := pattern
	if caseI { reStr = "(?i)" + reStr }
	re, err := regexp.Compile(reStr)
	if err != nil { fmt.Fprintln(os.Stderr, "kwsearch:", err); os.Exit(1) }

	const (
		matchColor = "\033[1;31m"
		ctxColor   = "\033[2m"
		reset      = "\033[0m"
	)

	search := func(r *os.File, name string) int {
		sc := bufio.NewScanner(r)
		var lines []string
		for sc.Scan() { lines = append(lines, sc.Text()) }
		matches := 0
		printed := map[int]bool{}
		for i, line := range lines {
			if !re.MatchString(line) { continue }
			matches++
			if countOnly { continue }
			start := i - before; if start < 0 { start = 0 }
			end := i + after; if end >= len(lines) { end = len(lines) - 1 }
			for j := start; j <= end; j++ {
				if printed[j] { continue }
				printed[j] = true
				prefix := "  "
				lineStr := lines[j]
				if j == i {
					prefix = "> "
					lineStr = re.ReplaceAllStringFunc(lineStr, func(m string) string {
						return matchColor + m + reset
					})
				}
				if showNums {
					fmt.Printf("%s%s:%d: %s%s\n", "", name, j+1, prefix, lineStr)
				} else {
					fmt.Printf("%s%s\n", prefix, lineStr)
				}
			}
			fmt.Println("--")
		}
		return matches
	}

	total := 0
	if len(files) == 0 {
		total = search(os.Stdin, "<stdin>")
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil { fmt.Fprintln(os.Stderr, err); continue }
			n := search(fh, f); total += n; fh.Close()
		}
	}
	if countOnly { fmt.Printf("Matches: %d\n", total) }
}
