package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func init() { register("wc", runWc) }

func runWc() {
	args := os.Args[1:]
	doLines, doWords, doBytes, doChars, doMaxLine := false, false, false, false, false
	files := []string{}

	for _, a := range args {
		switch a {
		case "-l", "--lines":
			doLines = true
		case "-w", "--words":
			doWords = true
		case "-c", "--bytes":
			doBytes = true
		case "-m", "--chars":
			doChars = true
		case "-L", "--max-line-length":
			doMaxLine = true
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}

	if !doLines && !doWords && !doBytes && !doChars && !doMaxLine {
		doLines, doWords, doBytes = true, true, true
	}

	type counts struct {
		lines, words, bytes, chars, maxLine int64
	}

	countReader := func(r io.Reader) counts {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		var c counts
		for sc.Scan() {
			line := sc.Text()
			c.lines++
			c.bytes += int64(len(line)) + 1
			c.chars += int64(len([]rune(line))) + 1
			c.words += int64(len(strings.Fields(line)))
			if int64(len(line)) > c.maxLine {
				c.maxLine = int64(len(line))
			}
		}
		return c
	}

	printCounts := func(c counts, name string) {
		if doLines {
			fmt.Printf("%8d", c.lines)
		}
		if doWords {
			fmt.Printf("%8d", c.words)
		}
		if doBytes {
			fmt.Printf("%8d", c.bytes)
		}
		if doChars {
			fmt.Printf("%8d", c.chars)
		}
		if doMaxLine {
			fmt.Printf("%8d", c.maxLine)
		}
		if name != "" {
			fmt.Printf(" %s", name)
		}
		fmt.Println()
	}

	if len(files) == 0 {
		c := countReader(os.Stdin)
		printCounts(c, "")
		return
	}

	var total counts
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wc: %s: %v\n", f, err)
			continue
		}
		c := countReader(fh)
		fh.Close()
		printCounts(c, f)
		total.lines += c.lines
		total.words += c.words
		total.bytes += c.bytes
		total.chars += c.chars
		if c.maxLine > total.maxLine {
			total.maxLine = c.maxLine
		}
	}
	if len(files) > 1 {
		printCounts(total, "total")
	}
}
