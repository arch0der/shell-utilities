package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() { register("rm", runRm) }

func runRm() {
	args := os.Args[1:]
	recursive := false
	force := false
	interactive := false
	verbose := false
	files := []string{}

	for _, a := range args {
		switch a {
		case "-r", "-R", "--recursive":
			recursive = true
		case "-f", "--force":
			force = true
		case "-i":
			interactive = true
		case "-v", "--verbose":
			verbose = true
		case "-rf", "-fr":
			recursive, force = true, true
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}

	confirm := func(path string) bool {
		if !interactive {
			return true
		}
		fmt.Fprintf(os.Stderr, "rm: remove '%s'? ", path)
		sc := bufio.NewScanner(os.Stdin)
		if sc.Scan() {
			resp := strings.ToLower(strings.TrimSpace(sc.Text()))
			return resp == "y" || resp == "yes"
		}
		return false
	}

	exitCode := 0
	var doRemove func(path string)
	doRemove = func(path string) {
		if !confirm(path) {
			return
		}
		info, err := os.Lstat(path)
		if err != nil {
			if !force {
				fmt.Fprintf(os.Stderr, "rm: %s: %v\n", path, err)
				exitCode = 1
			}
			return
		}
		if info.IsDir() {
			if !recursive {
				fmt.Fprintf(os.Stderr, "rm: cannot remove '%s': Is a directory\n", path)
				exitCode = 1
				return
			}
			entries, _ := os.ReadDir(path)
			for _, e := range entries {
				doRemove(filepath.Join(path, e.Name()))
			}
			if err := os.Remove(path); err != nil && !force {
				fmt.Fprintf(os.Stderr, "rm: %s: %v\n", path, err)
				exitCode = 1
			} else if verbose {
				fmt.Printf("removed directory '%s'\n", path)
			}
			return
		}
		if err := os.Remove(path); err != nil && !force {
			fmt.Fprintf(os.Stderr, "rm: %s: %v\n", path, err)
			exitCode = 1
		} else if verbose {
			fmt.Printf("removed '%s'\n", path)
		}
	}

	if len(files) == 0 && !force {
		fmt.Fprintln(os.Stderr, "rm: missing operand")
		os.Exit(1)
	}

	for _, f := range files {
		doRemove(f)
	}
	os.Exit(exitCode)
}
