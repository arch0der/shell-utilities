package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func init() { register("dirname", runDirname) }

func runDirname() {
	args := os.Args[1:]
	zero := false
	paths := []string{}
	for _, a := range args {
		if a == "-z" || a == "--zero" {
			zero = true
		} else {
			paths = append(paths, a)
		}
	}
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "dirname: missing operand")
		os.Exit(1)
	}
	sep := "\n"
	if zero {
		sep = "\x00"
	}
	for i, p := range paths {
		d := filepath.Dir(p)
		if i < len(paths)-1 || zero {
			fmt.Print(d + sep)
		} else {
			fmt.Println(d)
		}
	}
}
