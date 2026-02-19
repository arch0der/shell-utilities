// fold - Wrap lines to a given width
// Usage: fold [-w width] [-s] [-b] [file...]
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
	width    = flag.Int("w", 80, "Maximum line width")
	wordWrap = flag.Bool("s", false, "Break at spaces (word wrap)")
	bytes    = flag.Bool("b", false, "Count bytes instead of columns")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: fold [-w width] [-s] [-b] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	process := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			foldLine(sc.Text(), *width, *wordWrap)
		}
	}

	if flag.NArg() == 0 {
		process(os.Stdin)
		return
	}
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "fold:", err)
			continue
		}
		process(f)
		f.Close()
	}
}

func foldLine(line string, width int, wordWrap bool) {
	runes := []rune(line)
	if len(runes) == 0 {
		fmt.Println()
		return
	}

	for len(runes) > 0 {
		if len(runes) <= width {
			fmt.Println(string(runes))
			return
		}

		cut := width
		if wordWrap {
			// Find last space within width
			lastSpace := -1
			for i := 0; i < width && i < len(runes); i++ {
				if runes[i] == ' ' {
					lastSpace = i
				}
			}
			if lastSpace > 0 {
				cut = lastSpace + 1
			}
		}

		fmt.Println(strings.TrimRight(string(runes[:cut]), " "))
		runes = runes[cut:]
	}
}
