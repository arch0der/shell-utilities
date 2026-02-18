// tee - Read stdin and write to stdout and files simultaneously
// Usage: tee [-a] [file...]
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var appendMode = flag.Bool("a", false, "Append to files instead of overwriting")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: tee [-a] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	writers := []io.Writer{os.Stdout}

	flag := os.O_CREATE | os.O_WRONLY
	if *appendMode {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	for _, path := range os.Args[1:] {
		if path == "-a" {
			continue
		}
		f, err := os.OpenFile(path, flag, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tee:", err)
			os.Exit(1)
		}
		defer f.Close()
		writers = append(writers, f)
	}

	mw := io.MultiWriter(writers...)
	if _, err := io.Copy(mw, os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, "tee:", err)
		os.Exit(1)
	}
}
