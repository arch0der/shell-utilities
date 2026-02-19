package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	zero := false
	names := []string{}
	for _, a := range args {
		if a == "-0" || a == "--null" {
			zero = true
		} else if !strings.HasPrefix(a, "-") {
			names = append(names, a)
		}
	}
	sep := "\n"
	if zero {
		sep = "\x00"
	}
	if len(names) == 0 {
		for _, e := range os.Environ() {
			fmt.Print(e + sep)
		}
		return
	}
	exitCode := 0
	for _, name := range names {
		val, ok := os.LookupEnv(name)
		if !ok {
			exitCode = 1
		} else {
			fmt.Print(val + sep)
		}
	}
	os.Exit(exitCode)
}
