// bytecount - count bytes, chars, words, and lines (wc with readable output)
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type stats struct{ lines, words, chars int; bytes int64 }

func count(r io.Reader) stats {
	var s stats
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		s.lines++
		s.bytes += int64(len(line)) + 1
		s.chars += len([]rune(line)) + 1
		s.words += len(strings.Fields(line))
	}
	return s
}

func print1(name string, s stats) {
	fmt.Printf("%-24s  lines:%-8d  words:%-8d  bytes:%-10d  chars:%d\n",
		name, s.lines, s.words, s.bytes, s.chars)
}

func main() {
	if len(os.Args) == 1 {
		print1("<stdin>", count(os.Stdin)); return
	}
	var total stats
	for _, f := range os.Args[1:] {
		fh, err := os.Open(f)
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		s := count(fh); fh.Close()
		print1(f, s)
		total.lines += s.lines; total.words += s.words
		total.bytes += s.bytes; total.chars += s.chars
	}
	if len(os.Args) > 2 { print1("TOTAL", total) }
}
