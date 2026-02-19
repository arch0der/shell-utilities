// trim - Trim whitespace (or custom chars) from lines.
//
// Usage:
//
//	trim [OPTIONS] [FILE...]
//
// Options:
//
//	-l        Trim left only
//	-r        Trim right only
//	-c CHARS  Custom characters to trim (default: whitespace)
//	-n        Remove blank lines after trimming
//
// Examples:
//
//	cat file.txt | trim
//	trim -r file.txt
//	trim -c "\"'" file.txt     # trim quotes
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	left := flag.Bool("l", false, "trim left only")
	right := flag.Bool("r", false, "trim right only")
	chars := flag.String("c", "", "custom chars to trim")
	skipBlank := flag.Bool("n", false, "remove blank lines")
	flag.Parse()

	cutset := *chars

	trimLine := func(s string) string {
		if cutset != "" {
			if *left {
				return strings.TrimLeft(s, cutset)
			} else if *right {
				return strings.TrimRight(s, cutset)
			}
			return strings.Trim(s, cutset)
		}
		if *left {
			return strings.TrimLeftFunc(s, func(r rune) bool { return r == ' ' || r == '\t' })
		} else if *right {
			return strings.TrimRightFunc(s, func(r rune) bool { return r == ' ' || r == '\t' })
		}
		return strings.TrimSpace(s)
	}

	files := flag.Args()
	var readers []*os.File
	if len(files) == 0 {
		readers = []*os.File{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "trim: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for _, r := range readers {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			line := trimLine(sc.Text())
			if *skipBlank && line == "" {
				continue
			}
			fmt.Fprintln(w, line)
		}
	}
}
