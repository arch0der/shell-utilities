// head - Print first N lines of a file
// Usage: head [-n lines] [file...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var n = flag.Int("n", 10, "Number of lines to print")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: head [-n lines] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		head(os.Stdin, "")
		return
	}
	multi := flag.NArg() > 1
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "head:", err)
			continue
		}
		if multi {
			fmt.Printf("==> %s <==\n", path)
		}
		head(f, path)
		f.Close()
	}
}

func head(r io.Reader, _ string) {
	scanner := bufio.NewScanner(r)
	for i := 0; i < *n && scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
	}
}
