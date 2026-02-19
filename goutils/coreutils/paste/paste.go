package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	delim := "\t"
	serial := false
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-d" && i+1 < len(args):
			i++
			delim = args[i]
		case strings.HasPrefix(a, "-d"):
			delim = a[2:]
		case a == "-s" || a == "--serial":
			serial = true
		case a == "-z" || a == "--zero-terminated":
			delim = "\x00"
		case !strings.HasPrefix(a, "-") || a == "-":
			files = append(files, a)
		}
	}

	openReader := func(f string) io.Reader {
		if f == "-" {
			return os.Stdin
		}
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "paste: %s: %v\n", f, err)
			return strings.NewReader("")
		}
		return fh
	}

	if serial {
		for _, f := range files {
			r := openReader(f)
			sc := bufio.NewScanner(r)
			var parts []string
			for sc.Scan() {
				parts = append(parts, sc.Text())
			}
			fmt.Println(strings.Join(parts, delim))
		}
		return
	}

	scanners := make([]*bufio.Scanner, len(files))
	for i, f := range files {
		scanners[i] = bufio.NewScanner(openReader(f))
	}

	for {
		var parts []string
		any := false
		for _, sc := range scanners {
			if sc.Scan() {
				parts = append(parts, sc.Text())
				any = true
			} else {
				parts = append(parts, "")
			}
		}
		if !any {
			break
		}
		fmt.Println(strings.Join(parts, delim))
	}
}
