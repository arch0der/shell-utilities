// expand - Convert tabs to spaces
// Usage: expand [-t tabsize] [file...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var tabSize = flag.Int("t", 8, "Tab stop size")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: expand [-t tabsize] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	process := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			fmt.Println(expandTabs(sc.Text(), *tabSize))
		}
	}

	if flag.NArg() == 0 {
		process(os.Stdin)
		return
	}
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "expand:", err)
			continue
		}
		process(f)
		f.Close()
	}
}

func expandTabs(line string, size int) string {
	var sb strings.Builder
	col := 0
	for _, c := range line {
		if c == '\t' {
			spaces := size - (col % size)
			sb.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			sb.WriteRune(c)
			col++
		}
	}
	return sb.String()
}
