// wc - Count lines, words, and characters
// Usage: wc [-l] [-w] [-c] [file...]
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
	lines = flag.Bool("l", false, "Count lines")
	words = flag.Bool("w", false, "Count words")
	chars = flag.Bool("c", false, "Count characters")
)

type counts struct{ l, w, c int }

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: wc [-l] [-w] [-c] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	// If no flags, show all
	if !*lines && !*words && !*chars {
		*lines, *words, *chars = true, true, true
	}

	if flag.NArg() == 0 {
		c := count(os.Stdin)
		printCounts(c, "")
		return
	}

	total := counts{}
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wc:", err)
			continue
		}
		c := count(f)
		f.Close()
		printCounts(c, path)
		total.l += c.l
		total.w += c.w
		total.c += c.c
	}
	if flag.NArg() > 1 {
		printCounts(total, "total")
	}
}

func count(r io.Reader) counts {
	var c counts
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		c.l++
		line := scanner.Text()
		c.c += len(line) + 1
		c.w += len(strings.Fields(line))
	}
	return c
}

func printCounts(c counts, name string) {
	if *lines {
		fmt.Printf("%8d", c.l)
	}
	if *words {
		fmt.Printf("%8d", c.w)
	}
	if *chars {
		fmt.Printf("%8d", c.c)
	}
	if name != "" {
		fmt.Printf(" %s", name)
	}
	fmt.Println()
}
