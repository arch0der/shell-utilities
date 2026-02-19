// split - Split a file into pieces
// Usage: split [-l lines] [-b bytes] [-d] [file] [prefix]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	lines  = flag.Int("l", 1000, "Lines per output file")
	bytes  = flag.Int64("b", 0, "Bytes per output file (overrides -l)")
	numeric = flag.Bool("d", false, "Use numeric suffixes instead of alphabetic")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: split [-l lines] [-b bytes] [-d] [file] [prefix]")
		flag.PrintDefaults()
	}
	flag.Parse()

	prefix := "x"
	var r io.Reader = os.Stdin

	args := flag.Args()
	if len(args) > 0 && args[0] != "-" {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "split:", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}
	if len(args) > 1 {
		prefix = args[1]
	}

	if *bytes > 0 {
		splitBytes(r, prefix, *bytes)
	} else {
		splitLines(r, prefix, *lines)
	}
}

func suffix(n int, numeric bool) string {
	if numeric {
		return fmt.Sprintf("%02d", n)
	}
	a := n / 26
	b := n % 26
	if a > 25 {
		return fmt.Sprintf("zz%02d", n)
	}
	return fmt.Sprintf("%c%c", 'a'+a, 'a'+b)
}

func splitLines(r io.Reader, prefix string, n int) {
	sc := bufio.NewScanner(r)
	fileNum := 0
	count := 0
	var f *os.File

	for sc.Scan() {
		if count%n == 0 {
			if f != nil {
				f.Close()
			}
			name := prefix + suffix(fileNum, *numeric)
			var err error
			f, err = os.Create(name)
			if err != nil {
				fmt.Fprintln(os.Stderr, "split:", err)
				os.Exit(1)
			}
			fileNum++
		}
		fmt.Fprintln(f, sc.Text())
		count++
	}
	if f != nil {
		f.Close()
	}
}

func splitBytes(r io.Reader, prefix string, n int64) {
	buf := make([]byte, n)
	fileNum := 0
	for {
		nr, err := io.ReadFull(r, buf)
		if nr > 0 {
			name := prefix + suffix(fileNum, *numeric)
			os.WriteFile(name, buf[:nr], 0644)
			fileNum++
		}
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "split:", err)
			os.Exit(1)
		}
	}
}
