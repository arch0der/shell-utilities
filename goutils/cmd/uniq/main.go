// uniq - Remove or report duplicate lines
// Usage: uniq [-c] [-d] [-u] [file]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	count    = flag.Bool("c", false, "Prefix count of occurrences")
	dupsOnly = flag.Bool("d", false, "Only print duplicate lines")
	uniqOnly = flag.Bool("u", false, "Only print non-duplicate lines")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: uniq [-c] [-d] [-u] [file]")
		flag.PrintDefaults()
	}
	flag.Parse()

	var r io.Reader = os.Stdin
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "uniq:", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}

	scanner := bufio.NewScanner(r)
	var prev string
	cnt := 0
	first := true

	flush := func() {
		if first {
			return
		}
		if *dupsOnly && cnt == 1 {
			return
		}
		if *uniqOnly && cnt > 1 {
			return
		}
		if *count {
			fmt.Printf("%7d %s\n", cnt, prev)
		} else {
			fmt.Println(prev)
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		if first {
			prev = line
			cnt = 1
			first = false
			continue
		}
		if line == prev {
			cnt++
		} else {
			flush()
			prev = line
			cnt = 1
		}
	}
	flush()
}
