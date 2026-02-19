package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func init() { register("mkdir", runMkdir) }

func runMkdir() {
	args := os.Args[1:]
	parents := false
	mode := os.FileMode(0755)
	verbose := false
	dirs := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-p" || a == "--parents":
			parents = true
		case a == "-v" || a == "--verbose":
			verbose = true
		case a == "-m" && i+1 < len(args):
			i++
			m, _ := strconv.ParseUint(args[i], 8, 32)
			mode = os.FileMode(m)
		case strings.HasPrefix(a, "-m"):
			m, _ := strconv.ParseUint(a[2:], 8, 32)
			mode = os.FileMode(m)
		case !strings.HasPrefix(a, "-"):
			dirs = append(dirs, a)
		}
	}

	exitCode := 0
	for _, d := range dirs {
		var err error
		if parents {
			err = os.MkdirAll(d, mode)
		} else {
			err = os.Mkdir(d, mode)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "mkdir: %s: %v\n", d, err)
			exitCode = 1
		} else if verbose {
			fmt.Printf("mkdir: created directory '%s'\n", d)
		}
	}
	os.Exit(exitCode)
}
