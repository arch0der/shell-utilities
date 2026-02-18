// fmt - Simple text formatter / word wrapper
// Usage: fmt [-w width] [file...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var width = flag.Int("w", 75, "Maximum line width")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: fmt [-w width] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		formatReader(os.Stdin)
		return
	}
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "fmt:", err)
			continue
		}
		formatReader(f)
		f.Close()
	}
}

func formatReader(r io.Reader) {
	scanner := bufio.NewScanner(r)
	var para []string

	flush := func() {
		if len(para) == 0 {
			return
		}
		text := strings.Join(para, " ")
		words := strings.Fields(text)
		line := ""
		for _, w := range words {
			if line == "" {
				line = w
			} else if len(line)+1+len(w) <= *width {
				line += " " + w
			} else {
				fmt.Println(line)
				line = w
			}
		}
		if line != "" {
			fmt.Println(line)
		}
		para = nil
	}

	for scanner.Scan() {
		text := scanner.Text()
		if strings.TrimSpace(text) == "" {
			flush()
			fmt.Println()
		} else {
			para = append(para, strings.TrimSpace(text))
		}
	}
	flush()
}
