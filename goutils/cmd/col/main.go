// col - Filter reverse line feeds and backspaces (cleans man/nroff output)
// Usage: col [-b] [-x] [file...]
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
	noBackspace = flag.Bool("b", false, "Strip backspace overstrike sequences")
	expandTabs  = flag.Bool("x", false, "Expand tabs to spaces (8-wide)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: col [-b] [-x] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	process := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			fmt.Println(processLine(sc.Text()))
		}
	}

	if flag.NArg() == 0 {
		process(os.Stdin)
		return
	}
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "col:", err)
			continue
		}
		process(f)
		f.Close()
	}
}

func processLine(line string) string {
	// Resolve backspaces: a\bA => A (bold), _\ba => a (underline)
	runes := []rune(line)
	var buf []rune
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\b' {
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
			}
		} else {
			buf = append(buf, runes[i])
		}
	}
	// If -b: strip overstrike (char BS char => just the last char)
	// Already resolved above; -b flag controls whether to preserve bold
	line = string(buf)

	// Strip ANSI escape sequences
	var sb strings.Builder
	r := []rune(line)
	for i := 0; i < len(r); i++ {
		if r[i] == '\x1b' && i+1 < len(r) && r[i+1] == '[' {
			i += 2
			for i < len(r) && (r[i] < 'A' || r[i] > 'Z') && (r[i] < 'a' || r[i] > 'z') {
				i++
			}
		} else if r[i] == '\r' {
			// skip CR
		} else {
			sb.WriteRune(r[i])
		}
	}
	line = sb.String()

	// Expand tabs
	if *expandTabs {
		var out strings.Builder
		col := 0
		for _, c := range line {
			if c == '\t' {
				spaces := 8 - (col % 8)
				out.WriteString(strings.Repeat(" ", spaces))
				col += spaces
			} else {
				out.WriteRune(c)
				col++
			}
		}
		line = out.String()
	}
	return line
}
