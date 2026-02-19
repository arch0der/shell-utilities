package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func init() { register("tee", runTee) }

func runTee() {
	args := os.Args[1:]
	appendMode := false
	ignoreInterrupt := false
	files := []string{}
	for _, a := range args {
		switch a {
		case "-a", "--append":
			appendMode = true
		case "-i", "--ignore-interrupts":
			ignoreInterrupt = true
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}
	_ = ignoreInterrupt

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if appendMode {
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}

	writers := []io.Writer{os.Stdout}
	for _, f := range files {
		fh, err := os.OpenFile(f, flag, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tee: %s: %v\n", f, err)
			continue
		}
		defer fh.Close()
		writers = append(writers, fh)
	}
	io.Copy(io.MultiWriter(writers...), os.Stdin)
}
