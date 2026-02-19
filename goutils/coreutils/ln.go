package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() { register("ln", runLn) }

func runLn() {
	args := os.Args[1:]
	symbolic := false
	force := false
	noDeref := false
	verbose := false
	relative := false
	files := []string{}

	for _, a := range args {
		switch a {
		case "-s", "--symbolic":
			symbolic = true
		case "-f", "--force":
			force = true
		case "-n", "--no-dereference":
			noDeref = true
		case "-v", "--verbose":
			verbose = true
		case "-r", "--relative":
			relative = true
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}
	_ = noDeref
	_ = relative

	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "ln: missing destination")
		os.Exit(1)
	}

	dest := files[len(files)-1]
	srcs := files[:len(files)-1]

	exitCode := 0
	for _, src := range srcs {
		dst := dest
		if info, err := os.Stat(dest); err == nil && info.IsDir() {
			dst = filepath.Join(dest, filepath.Base(src))
		}
		if force {
			os.Remove(dst)
		}
		var err error
		if symbolic {
			err = os.Symlink(src, dst)
		} else {
			err = os.Link(src, dst)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ln: %v\n", err)
			exitCode = 1
		} else if verbose {
			fmt.Printf("'%s' -> '%s'\n", src, dst)
		}
	}
	os.Exit(exitCode)
}
