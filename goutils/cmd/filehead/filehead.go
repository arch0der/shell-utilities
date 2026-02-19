// filehead - show the first N bytes of files (binary-friendly head)
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	n := int64(512)
	files := os.Args[1:]
	if len(files) > 0 {
		if v, err := strconv.ParseInt(files[0], 10, 64); err == nil { n = v; files = files[1:] }
	}
	if len(files) == 0 { files = []string{"-"} }

	for _, path := range files {
		if len(files) > 1 { fmt.Printf("==> %s <==\n", path) }
		var r io.Reader
		if path == "-" {
			r = os.Stdin
		} else {
			f, err := os.Open(path)
			if err != nil { fmt.Fprintln(os.Stderr, err); continue }
			defer f.Close()
			r = f
		}
		buf := make([]byte, n)
		got, err := io.ReadFull(r, buf)
		if err != nil && got == 0 { fmt.Fprintln(os.Stderr, err); continue }
		os.Stdout.Write(buf[:got])
		if len(files) > 1 { fmt.Println() }
	}
}
