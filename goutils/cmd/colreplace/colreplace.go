// colreplace - replace values in a specific column of delimited input
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
	fmt.Fprintln(os.Stderr, `usage: colreplace [options] <col> <pattern> <replacement>
  -d <delim>    field delimiter (default: tab)
  -H            first line is header; use name instead of number
  -e            pattern is a regexp (default: literal)
  Column is 1-based index or header name (with -H).`)
	os.Exit(1)
}

func main() {
	delim := "\t"
	useHeader := false
	useRegex := false
	args := os.Args[1:]
	rest := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d": i++; delim = args[i]
		case "-H": useHeader = true
		case "-e": useRegex = true
		default: rest = append(rest, args[i])
		}
	}
	if len(rest) < 3 { usage() }
	colSpec, pattern, replacement := rest[0], rest[1], rest[2]

	var re *regexp.Regexp
	if useRegex {
		var err error; re, err = regexp.Compile(pattern)
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	}

	colIdx := -1
	if n, err := strconv.Atoi(colSpec); err == nil { colIdx = n - 1 }

	sc := bufio.NewScanner(os.Stdin)
	first := true
	for sc.Scan() {
		line := sc.Text()
		fields := strings.Split(line, delim)
		if first && useHeader {
			first = false
			if colIdx < 0 {
				for i, h := range fields {
					if h == colSpec { colIdx = i; break }
				}
				if colIdx < 0 { fmt.Fprintf(os.Stderr, "colreplace: column %q not found\n", colSpec); os.Exit(1) }
			}
			fmt.Println(line); continue
		}
		first = false
		if colIdx >= 0 && colIdx < len(fields) {
			if useRegex {
				fields[colIdx] = re.ReplaceAllString(fields[colIdx], replacement)
			} else {
				fields[colIdx] = strings.ReplaceAll(fields[colIdx], pattern, replacement)
			}
		}
		fmt.Println(strings.Join(fields, delim))
	}
}
