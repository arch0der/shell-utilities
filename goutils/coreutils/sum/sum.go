package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	sysV := false
	files := []string{}
	for _, a := range args {
		if a == "-s" || a == "--sysv" {
			sysV = true
		} else if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}

	sumFile := func(r io.Reader, name string) {
		data, _ := io.ReadAll(r)
		var checksum uint32
		if sysV {
			for _, b := range data {
				checksum += uint32(b)
			}
			checksum = (checksum & 0xffff) + (checksum >> 16)
			checksum = (checksum & 0xffff) + (checksum >> 16)
		} else {
			// BSD sum
			for _, b := range data {
				checksum = (checksum >> 1) | (checksum << 15)
				checksum = (checksum + uint32(b)) & 0xffff
			}
		}
		blocks := (len(data) + 511) / 512
		if name != "" {
			fmt.Printf("%05d %5d %s\n", checksum, blocks, name)
		} else {
			fmt.Printf("%05d %5d\n", checksum, blocks)
		}
	}

	if len(files) == 0 {
		sumFile(os.Stdin, "")
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sum: %s: %v\n", f, err)
			continue
		}
		sumFile(fh, f)
		fh.Close()
	}
}
