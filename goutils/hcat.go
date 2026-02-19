// hcat - Display files side by side (horizontal cat).
//
// Usage:
//
//	hcat [OPTIONS] FILE [FILE...]
//
// Options:
//
//	-w N      Column width per file (default: auto from terminal)
//	-s SEP    Column separator (default: " │ ")
//	-n        Show line numbers
//	-t        Show filename headers
//	-p        Pad short files with empty lines to match longest
//
// Examples:
//
//	hcat file1.txt file2.txt
//	hcat -s " | " left.txt right.txt
//	hcat -t -w 40 original.txt modified.txt
//	diff <(cat a) <(cat b) | hcat -t a b     # side-by-side diff helper
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
	colWidth = flag.Int("w", 0, "column width")
	sep      = flag.String("s", " │ ", "separator")
	lineNums = flag.Bool("n", false, "line numbers")
	headers  = flag.Bool("t", false, "show headers")
	pad      = flag.Bool("p", true, "pad short files")
)

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

func padLine(s string, width int) string {
	n := utf8.RuneCountInString(s)
	if n >= width {
		runes := []rune(s)
		return string(runes[:width])
	}
	return s + strings.Repeat(" ", width-n)
}

func main() {
	flag.Parse()
	files := flag.Args()
	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "usage: hcat [OPTIONS] FILE FILE [FILE...]")
		os.Exit(1)
	}

	// Read all files
	var allLines [][]string
	maxLines := 0
	for _, f := range files {
		lines, err := readLines(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hcat: %v\n", err)
			os.Exit(1)
		}
		allLines = append(allLines, lines)
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}

	// Determine column width
	w := *colWidth
	if w == 0 {
		// Try terminal width, divide by num files
		termWidth := 120 // default
		numFiles := len(files)
		sepWidth := utf8.RuneCountInString(*sep)
		w = (termWidth - sepWidth*(numFiles-1)) / numFiles
		if w < 20 {
			w = 40
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	// Headers
	if *headers {
		parts := make([]string, len(files))
		for i, f := range files {
			parts[i] = padLine(f, w)
		}
		fmt.Fprintln(out, strings.Join(parts, *sep))
		// Divider
		divParts := make([]string, len(files))
		for i := range files {
			_ = i
			divParts[i] = strings.Repeat("─", w)
		}
		fmt.Fprintln(out, strings.Join(divParts, "─┼─"))
	}

	// Lines
	lineNumWidth := len(fmt.Sprintf("%d", maxLines))
	for i := 0; i < maxLines; i++ {
		parts := make([]string, len(allLines))
		for j, lines := range allLines {
			line := ""
			if i < len(lines) {
				line = lines[i]
			} else if !*pad {
				line = ""
			}
			parts[j] = padLine(line, w)
		}
		if *lineNums {
			fmt.Fprintf(out, "%*d  ", lineNumWidth, i+1)
		}
		fmt.Fprintln(out, strings.Join(parts, *sep))
	}
}
