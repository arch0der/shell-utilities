// unexpand - Convert leading spaces to tabs
// Usage: unexpand [-t tabsize] [-a] [file...]
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
	tabSize = flag.Int("t", 8, "Tab stop size")
	all     = flag.Bool("a", false, "Convert all whitespace, not just leading")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: unexpand [-t tabsize] [-a] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	process := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			fmt.Println(unexpandLine(sc.Text(), *tabSize, *all))
		}
	}

	if flag.NArg() == 0 {
		process(os.Stdin)
		return
	}
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unexpand:", err)
			continue
		}
		process(f)
		f.Close()
	}
}

func unexpandLine(line string, size int, all bool) string {
	// Convert leading spaces to tabs (or all spaces if -a)
	var sb strings.Builder
	col := 0
	i := 0
	runes := []rune(line)

	for i < len(runes) {
		c := runes[i]
		if c == ' ' || c == '\t' {
			// Count spaces
			spaces := 0
			j := i
			for j < len(runes) && (runes[j] == ' ' || runes[j] == '\t') {
				if runes[j] == '\t' {
					spaces += size - (spaces % size)
				} else {
					spaces++
				}
				j++
			}
			// Convert to tabs + remaining spaces
			tabs := spaces / size
			rem := spaces % size
			sb.WriteString(strings.Repeat("\t", tabs))
			sb.WriteString(strings.Repeat(" ", rem))
			col += spaces
			i = j
			if !all {
				// Only convert leading whitespace; rest verbatim
				for i < len(runes) {
					sb.WriteRune(runes[i])
					i++
				}
			}
		} else {
			sb.WriteRune(c)
			col++
			i++
		}
	}
	return sb.String()
}
