package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func init() { register("shred", runShred) }

func runShred() {
	args := os.Args[1:]
	iterations := 3
	verbose := false
	doRemove := false
	zero := false
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-n" && i+1 < len(args):
			i++
			iterations, _ = strconv.Atoi(args[i])
		case strings.HasPrefix(a, "-n"):
			iterations, _ = strconv.Atoi(a[2:])
		case a == "-v" || a == "--verbose":
			verbose = true
		case a == "-u":
			doRemove = true
		case a == "-z" || a == "--zero":
			zero = true
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}

	exitCode := 0
	for _, f := range files {
		fh, err := os.OpenFile(f, os.O_WRONLY, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "shred: %s: %v\n", f, err)
			exitCode = 1
			continue
		}
		info, _ := fh.Stat()
		size := info.Size()
		buf := make([]byte, 4096)
		for pass := 0; pass < iterations; pass++ {
			if verbose {
				fmt.Fprintf(os.Stderr, "shred: %s: pass %d/%d (random)\n", f, pass+1, iterations)
			}
			fh.Seek(0, 0)
			remaining := size
			for remaining > 0 {
				n := int64(len(buf))
				if n > remaining {
					n = remaining
					buf = buf[:n]
				}
				rand.Read(buf)
				fh.Write(buf)
				remaining -= n
			}
		}
		if zero {
			fh.Seek(0, 0)
			for i := range buf {
				buf[i] = 0
			}
			remaining := size
			for remaining > 0 {
				n := int64(len(buf))
				if n > remaining {
					n = remaining
				}
				fh.Write(buf[:n])
				remaining -= n
			}
		}
		fh.Close()
		if doRemove {
			os.Remove(f)
		}
	}
	os.Exit(exitCode)
}
