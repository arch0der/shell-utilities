// filetail - show the last N bytes of files (binary-friendly tail)
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
		if path == "-" {
			data, _ := io.ReadAll(os.Stdin)
			if int64(len(data)) > n { data = data[int64(len(data))-n:] }
			os.Stdout.Write(data)
		} else {
			f, err := os.Open(path)
			if err != nil { fmt.Fprintln(os.Stderr, err); continue }
			info, _ := f.Stat()
			size := info.Size()
			off := size - n
			if off < 0 { off = 0 }
			f.Seek(off, io.SeekStart)
			io.Copy(os.Stdout, f)
			f.Close()
		}
		if len(files) > 1 { fmt.Println() }
	}
}
