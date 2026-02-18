// nl - Number lines of files
// Usage: nl [-b a|t|n] [-n ln|rn|rz] [-w width] [-v start] [-i incr] [file...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	bodyStyle = flag.String("b", "t", "Body numbering: a=all, t=non-empty, n=none")
	numFormat = flag.String("n", "rn", "Number format: ln=left, rn=right, rz=right+zeros")
	numWidth  = flag.Int("w", 6, "Number field width")
	startNum  = flag.Int("v", 1, "Starting line number")
	increment = flag.Int("i", 1, "Line number increment")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: nl [-b a|t|n] [-n ln|rn|rz] [-w width] [-v start] [-i incr] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		numberLines(os.Stdin)
		return
	}
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "nl:", err)
			continue
		}
		numberLines(f)
		f.Close()
	}
}

func numberLines(r io.Reader) {
	scanner := bufio.NewScanner(r)
	lineNo := *startNum
	for scanner.Scan() {
		line := scanner.Text()
		shouldNumber := false
		switch *bodyStyle {
		case "a":
			shouldNumber = true
		case "t":
			shouldNumber = strings.TrimSpace(line) != ""
		case "n":
			shouldNumber = false
		}
		if shouldNumber {
			numStr := formatNum(lineNo, *numWidth, *numFormat)
			fmt.Printf("%s\t%s\n", numStr, line)
			lineNo += *increment
		} else {
			fmt.Printf("%s\t%s\n", strings.Repeat(" ", *numWidth), line)
		}
	}
}

func formatNum(n, w int, format string) string {
	switch format {
	case "ln":
		return fmt.Sprintf("%-*d", w, n)
	case "rz":
		return fmt.Sprintf("%0*d", w, n)
	default: // rn
		return fmt.Sprintf("%*d", w, n)
	}
}
