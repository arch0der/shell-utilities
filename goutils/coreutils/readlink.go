package main

import (
	"fmt"
	"os"
	"strings"
)

func init() { register("readlink", runReadlink) }

func runReadlink() {
	args := os.Args[1:]
	canonicalize := false
	canonicalizeMissing := false
	noNewline := false
	quiet := false
	files := []string{}

	for _, a := range args {
		switch a {
		case "-f", "--canonicalize":
			canonicalize = true
		case "-m", "--canonicalize-missing":
			canonicalizeMissing = true
		case "-n", "--no-newline":
			noNewline = true
		case "-q", "--quiet", "-s", "--silent":
			quiet = true
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}
	_ = quiet

	exitCode := 0
	for _, f := range files {
		var result string
		var err error
		if canonicalize || canonicalizeMissing {
			result, err = realPath(f)
		} else {
			result, err = os.Readlink(f)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "readlink: %s: %v\n", f, err)
			exitCode = 1
			continue
		}
		if noNewline {
			fmt.Print(result)
		} else {
			fmt.Println(result)
		}
	}
	os.Exit(exitCode)
}
