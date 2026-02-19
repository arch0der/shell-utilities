// paste - Merge lines of files side by side
// Usage: paste [-d delimiter] [-s] file...
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
	delim   = flag.String("d", "\t", "Delimiter between columns")
	serial  = flag.Bool("s", false, "Paste each file serially rather than in parallel")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: paste [-d delim] [-s] file...")
		flag.PrintDefaults()
	}
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		files = []string{"-"}
	}

	if *serial {
		pasteSerial(files)
	} else {
		pasteParallel(files)
	}
}

func openFile(path string) io.Reader {
	if path == "-" {
		return os.Stdin
	}
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "paste:", err)
		os.Exit(1)
	}
	return f
}

func pasteParallel(files []string) {
	scanners := make([]*bufio.Scanner, len(files))
	for i, f := range files {
		scanners[i] = bufio.NewScanner(openFile(f))
	}
	for {
		active := false
		parts := make([]string, len(scanners))
		for i, sc := range scanners {
			if sc.Scan() {
				parts[i] = sc.Text()
				active = true
			}
		}
		if !active {
			break
		}
		fmt.Println(strings.Join(parts, *delim))
	}
}

func pasteSerial(files []string) {
	for _, f := range files {
		sc := bufio.NewScanner(openFile(f))
		var parts []string
		for sc.Scan() {
			parts = append(parts, sc.Text())
		}
		fmt.Println(strings.Join(parts, *delim))
	}
}
