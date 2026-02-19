package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	parents := false
	verbose := false
	dirs := []string{}

	for _, a := range args {
		switch a {
		case "-p", "--parents":
			parents = true
		case "-v", "--verbose":
			verbose = true
		default:
			if !strings.HasPrefix(a, "-") {
				dirs = append(dirs, a)
			}
		}
	}

	exitCode := 0
	for _, d := range dirs {
		if err := os.Remove(d); err != nil {
			fmt.Fprintf(os.Stderr, "rmdir: %s: %v\n", d, err)
			exitCode = 1
			continue
		}
		if verbose {
			fmt.Printf("rmdir: removing directory '%s'\n", d)
		}
		if parents {
			p := filepath.Dir(d)
			for p != "." && p != "/" {
				if err := os.Remove(p); err != nil {
					break
				}
				if verbose {
					fmt.Printf("rmdir: removing directory '%s'\n", p)
				}
				p = filepath.Dir(p)
			}
		}
	}
	os.Exit(exitCode)
}
