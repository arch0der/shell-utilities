// wordcount - Count lines, words, characters, and bytes in files.
//
// Usage:
//
//	wordcount [OPTIONS] [FILE...]
//	cat file | wordcount
//
// Options:
//
//	-l        Lines only
//	-w        Words only
//	-c        Characters only
//	-b        Bytes only
//	-j        JSON output
//	-H        Human-readable numbers
//
// Examples:
//
//	wordcount file.txt
//	cat access.log | wordcount -l
//	wordcount -j *.go
//	wordcount -H large_file.txt
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

var (
	linesOnly = flag.Bool("l", false, "lines only")
	wordsOnly = flag.Bool("w", false, "words only")
	charsOnly = flag.Bool("c", false, "chars only")
	bytesOnly = flag.Bool("b", false, "bytes only")
	asJSON    = flag.Bool("j", false, "JSON output")
	human     = flag.Bool("H", false, "human-readable")
)

type Counts struct {
	File  string `json:"file"`
	Lines int64  `json:"lines"`
	Words int64  `json:"words"`
	Chars int64  `json:"chars"`
	Bytes int64  `json:"bytes"`
}

func humanize(n int64) string {
	switch {
	case n >= 1_000_000_000:
		return fmt.Sprintf("%.1fG", float64(n)/1e9)
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1e3)
	}
	return fmt.Sprintf("%d", n)
}

func count(r *os.File, name string) Counts {
	c := Counts{File: name}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		c.Lines++
		c.Words += int64(len(strings.Fields(line)))
		c.Chars += int64(utf8.RuneCountInString(line)) + 1 // +1 for newline
		c.Bytes += int64(len(line)) + 1
	}
	return c
}

func printCounts(c Counts) {
	n := func(v int64) string {
		if *human {
			return humanize(v)
		}
		return fmt.Sprintf("%8d", v)
	}

	any := *linesOnly || *wordsOnly || *charsOnly || *bytesOnly
	if !any {
		fmt.Printf("%s %s %s %s  %s\n", n(c.Lines), n(c.Words), n(c.Chars), n(c.Bytes), c.File)
		return
	}
	parts := []string{}
	if *linesOnly {
		parts = append(parts, n(c.Lines))
	}
	if *wordsOnly {
		parts = append(parts, n(c.Words))
	}
	if *charsOnly {
		parts = append(parts, n(c.Chars))
	}
	if *bytesOnly {
		parts = append(parts, n(c.Bytes))
	}
	fmt.Printf("%s  %s\n", strings.Join(parts, " "), c.File)
}

func main() {
	flag.Parse()
	files := flag.Args()

	var all []Counts

	if len(files) == 0 {
		c := count(os.Stdin, "stdin")
		all = append(all, c)
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "wordcount: %v\n", err)
				continue
			}
			c := count(fh, f)
			fh.Close()
			all = append(all, c)
		}
	}

	if *asJSON {
		b, _ := json.MarshalIndent(all, "", "  ")
		fmt.Println(string(b))
		return
	}

	// Print header if no specific flag
	if !*linesOnly && !*wordsOnly && !*charsOnly && !*bytesOnly {
		fmt.Printf("%8s %8s %8s %8s  %s\n", "lines", "words", "chars", "bytes", "file")
	}

	var total Counts
	total.File = "total"
	for _, c := range all {
		printCounts(c)
		total.Lines += c.Lines
		total.Words += c.Words
		total.Chars += c.Chars
		total.Bytes += c.Bytes
	}
	if len(all) > 1 {
		printCounts(total)
	}
}
