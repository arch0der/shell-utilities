// truncate - Truncate lines to a maximum width.
//
// Usage:
//
//	truncate [OPTIONS] [FILE...]
//	cat file | truncate -n 80
//
// Options:
//
//	-n N      Maximum line length (default: 80)
//	-e STR    Ellipsis string appended when truncated (default: ...)
//	-b        Truncate from the beginning (keep end)
//	-w        Truncate at word boundary
//
// Examples:
//
//	cat long.txt | truncate -n 60
//	cat log.txt  | truncate -n 120 -e "…"
//	truncate -n 40 -w file.txt       # word-boundary truncation
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

var (
	maxLen   = flag.Int("n", 80, "max line length")
	ellipsis = flag.String("e", "...", "ellipsis")
	fromEnd  = flag.Bool("b", false, "truncate from beginning")
	word     = flag.Bool("w", false, "word boundary")
)

func truncateLine(s string) string {
	runes := []rune(s)
	l := len(runes)
	el := utf8.RuneCountInString(*ellipsis)
	max := *maxLen

	if l <= max {
		return s
	}

	if *fromEnd {
		// keep end
		keep := max - el
		if keep < 0 {
			keep = 0
		}
		return *ellipsis + string(runes[l-keep:])
	}

	keep := max - el
	if keep < 0 {
		keep = 0
	}

	if *word {
		// find last space at or before keep
		truncAt := keep
		for truncAt > 0 && runes[truncAt-1] != ' ' {
			truncAt--
		}
		if truncAt == 0 {
			truncAt = keep // no word boundary found, hard cut
		}
		return string(runes[:truncAt]) + *ellipsis
	}

	return string(runes[:keep]) + *ellipsis
}

func main() {
	flag.Parse()
	files := flag.Args()

	var readers []*os.File
	if len(files) == 0 {
		readers = []*os.File{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "truncate: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	_ = strings.TrimSpace

	for _, r := range readers {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1024*1024), 1024*1024)
		for sc.Scan() {
			fmt.Fprintln(w, truncateLine(sc.Text()))
		}
	}
}
