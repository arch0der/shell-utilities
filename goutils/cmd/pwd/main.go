// pwd - Print working directory
// Usage: pwd [-P]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var physical = flag.Bool("P", false, "Print physical path (resolve symlinks)")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: pwd [-P]")
		flag.PrintDefaults()
	}
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pwd:", err)
		os.Exit(1)
	}

	if *physical {
		dir, err = filepath.EvalSymlinks(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "pwd:", err)
			os.Exit(1)
		}
	} else {
		if env := os.Getenv("PWD"); env != "" {
			dir = env
		}
	}

	fmt.Println(dir)
}
