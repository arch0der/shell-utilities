// dedup - Remove duplicate lines, preserving original order.
//
// Usage:
//
//	dedup [OPTIONS] [FILE...]
//
// Options:
//
//	-i        Case-insensitive comparison
//	-c        Prefix lines with occurrence count
//	-d        Print only duplicate lines (seen >1 times)
//	-u        Print only unique lines (seen exactly once)
//	-f N      Skip first N fields when comparing
//
// Examples:
//
//	cat file.txt | dedup
//	dedup -i -c file.txt
//	dedup -d access.log
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	caseInsensitive := flag.Bool("i", false, "case-insensitive")
	showCount := flag.Bool("c", false, "prefix with count")
	onlyDups := flag.Bool("d", false, "only duplicates")
	onlyUniq := flag.Bool("u", false, "only unique")
	skipFields := flag.Int("f", 0, "skip first N fields")
	flag.Parse()

	files := flag.Args()
	var readers []*os.File
	if len(files) == 0 {
		readers = []*os.File{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dedup: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}

	type entry struct {
		line  string
		count int
	}

	seen := make(map[string]*entry)
	var order []string

	for _, r := range readers {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			line := sc.Text()
			key := line
			if *skipFields > 0 {
				parts := strings.Fields(line)
				if *skipFields < len(parts) {
					key = strings.Join(parts[*skipFields:], " ")
				} else {
					key = ""
				}
			}
			if *caseInsensitive {
				key = strings.ToLower(key)
			}
			if e, ok := seen[key]; ok {
				e.count++
			} else {
				seen[key] = &entry{line: line, count: 1}
				order = append(order, key)
			}
		}
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for _, key := range order {
		e := seen[key]
		if *onlyDups && e.count < 2 {
			continue
		}
		if *onlyUniq && e.count > 1 {
			continue
		}
		if *showCount {
			fmt.Fprintf(w, "%7d %s\n", e.count, e.line)
		} else {
			fmt.Fprintln(w, e.line)
		}
	}
}
