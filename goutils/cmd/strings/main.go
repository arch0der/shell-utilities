// strings - Find printable strings in binary files
// Usage: strings [-n minlen] [-o] [file...]
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	minLen  = flag.Int("n", 4, "Minimum string length")
	offsets = flag.Bool("o", false, "Print file offsets")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: strings [-n minlen] [-o] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		findStrings(os.Stdin, "")
		return
	}
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "strings:", err)
			continue
		}
		findStrings(f, path)
		f.Close()
	}
}

func isPrintable(b byte) bool {
	return b >= 0x20 && b < 0x7f
}

func findStrings(r io.Reader, filename string) {
	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "strings:", err)
		return
	}
	var cur strings.Builder
	start := 0
	for i, b := range data {
		if isPrintable(b) {
			if cur.Len() == 0 {
				start = i
			}
			cur.WriteByte(b)
		} else {
			if cur.Len() >= *minLen {
				if filename != "" {
					fmt.Printf("%s: ", filename)
				}
				if *offsets {
					fmt.Printf("%d: ", start)
				}
				fmt.Println(cur.String())
			}
			cur.Reset()
		}
	}
	if cur.Len() >= *minLen {
		if filename != "" {
			fmt.Printf("%s: ", filename)
		}
		if *offsets {
			fmt.Printf("%d: ", start)
		}
		fmt.Println(cur.String())
	}
}
