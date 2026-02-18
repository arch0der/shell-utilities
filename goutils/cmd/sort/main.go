// sort - Sort lines of text files
// Usage: sort [-r] [-n] [-u] [-k field] [file...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	reverse = flag.Bool("r", false, "Reverse sort order")
	numeric = flag.Bool("n", false, "Numeric sort")
	unique  = flag.Bool("u", false, "Remove duplicate lines")
	field   = flag.Int("k", 0, "Sort by field number (1-indexed, 0=whole line)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: sort [-r] [-n] [-u] [-k field] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	var lines []string
	if flag.NArg() == 0 {
		lines = readLines(os.Stdin)
	} else {
		for _, path := range flag.Args() {
			f, err := os.Open(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "sort:", err)
				continue
			}
			lines = append(lines, readLines(f)...)
			f.Close()
		}
	}

	sort.SliceStable(lines, func(i, j int) bool {
		a, b := getKey(lines[i]), getKey(lines[j])
		if *numeric {
			na, _ := strconv.ParseFloat(a, 64)
			nb, _ := strconv.ParseFloat(b, 64)
			if *reverse {
				return na > nb
			}
			return na < nb
		}
		if *reverse {
			return a > b
		}
		return a < b
	})

	seen := map[string]bool{}
	for _, line := range lines {
		if *unique {
			if seen[line] {
				continue
			}
			seen[line] = true
		}
		fmt.Println(line)
	}
}

func getKey(line string) string {
	if *field == 0 {
		return line
	}
	parts := strings.Fields(line)
	if *field <= len(parts) {
		return parts[*field-1]
	}
	return ""
}

func readLines(r io.Reader) []string {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
