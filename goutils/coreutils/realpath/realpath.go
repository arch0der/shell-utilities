package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func realPath(p string) (string, error) {
	return filepath.EvalSymlinks(p)
}

func main() {
	args := os.Args[1:]
	quiet := false
	relative := false
	logicalOnly := false
	missing := false
	files := []string{}

	for _, a := range args {
		switch a {
		case "-q", "--quiet":
			quiet = true
		case "-s", "--strip", "--no-symlinks":
			logicalOnly = true
		case "--relative-to":
			relative = true
		case "-m", "--canonicalize-missing":
			missing = true
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}
	_ = relative
	_ = logicalOnly
	_ = missing

	exitCode := 0
	for _, f := range files {
		var result string
		var err error
		result, err = filepath.Abs(f)
		if err != nil {
			if !quiet {
				fmt.Fprintf(os.Stderr, "realpath: %s: %v\n", f, err)
			}
			exitCode = 1
			continue
		}
		// Try to resolve symlinks
		if resolved, err2 := filepath.EvalSymlinks(result); err2 == nil {
			result = resolved
		}
		fmt.Println(result)
	}
	os.Exit(exitCode)
}
