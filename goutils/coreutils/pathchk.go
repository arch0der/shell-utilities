package main

import (
	"fmt"
	"os"
	"strings"
)

func init() { register("pathchk", runPathchk) }

func runPathchk() {
	args := os.Args[1:]
	posix := false
	portability := false
	paths := []string{}

	for _, a := range args {
		switch a {
		case "-p":
			posix = true
		case "-P":
			portability = true
		default:
			if !strings.HasPrefix(a, "-") {
				paths = append(paths, a)
			}
		}
	}
	_ = posix
	_ = portability

	exitCode := 0
	for _, p := range paths {
		if len(p) == 0 {
			fmt.Fprintln(os.Stderr, "pathchk: empty path")
			exitCode = 1
			continue
		}
		if len(p) > 4096 {
			fmt.Fprintf(os.Stderr, "pathchk: '%s': path too long\n", p)
			exitCode = 1
		}
		for _, component := range strings.Split(p, "/") {
			if len(component) > 255 {
				fmt.Fprintf(os.Stderr, "pathchk: '%s': component too long\n", p)
				exitCode = 1
			}
		}
	}
	os.Exit(exitCode)
}
