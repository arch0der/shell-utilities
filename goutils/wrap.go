// wrap - Word-wrap text at a specified column width.
//
// Usage:
//
//	wrap [OPTIONS] [FILE...]
//	echo "long text..." | wrap -n 72
//
// Options:
//
//	-n N      Wrap at column N (default: 80)
//	-i STR    Indent each line with STR
//	-p STR    Prefix first line with STR (like paragraph indent)
//	-s        Strict: don't break long words (leave them as-is)
//
// Examples:
//
//	cat README.txt | wrap -n 72
//	echo "..." | wrap -n 60 -i "  "        # 2-space indent
//	wrap -n 72 -p "> " email.txt           # blockquote style
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
	width  = flag.Int("n", 80, "wrap width")
	indent = flag.String("i", "", "indent")
	prefix = flag.String("p", "", "first line prefix")
	strict = flag.Bool("s", false, "don't break long words")
)

func wrapLine(line string) []string {
	if line == "" {
		return []string{""}
	}
	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{""}
	}

	indentW := utf8.RuneCountInString(*indent)
	prefixW := utf8.RuneCountInString(*prefix)
	maxW := *width

	var lines []string
	var cur strings.Builder
	curLen := 0
	firstLine := true

	flush := func() {
		if firstLine {
			lines = append(lines, *prefix+cur.String())
			firstLine = false
		} else {
			lines = append(lines, *indent+cur.String())
		}
		cur.Reset()
		curLen = 0
	}

	for _, word := range words {
		wLen := utf8.RuneCountInString(word)
		extra := indentW
		if firstLine {
			extra = prefixW
		}
		available := maxW - extra

		if curLen == 0 {
			if wLen > available && !*strict {
				cur.WriteString(word)
				curLen = wLen
				flush()
			} else {
				cur.WriteString(word)
				curLen = wLen
			}
		} else if curLen+1+wLen <= available {
			cur.WriteString(" " + word)
			curLen += 1 + wLen
		} else {
			flush()
			if wLen > maxW-indentW && !*strict {
				cur.WriteString(word)
				curLen = wLen
				flush()
			} else {
				cur.WriteString(word)
				curLen = wLen
			}
		}
	}
	if curLen > 0 {
		flush()
	}
	return lines
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
				fmt.Fprintf(os.Stderr, "wrap: %v\n", err)
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
		sc.Buffer(make([]byte, 1024*1024), 1024*1024)
		for sc.Scan() {
			for _, l := range wrapLine(sc.Text()) {
				fmt.Fprintln(w, l)
			}
		}
	}
}
