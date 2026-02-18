// rm - Remove files or directories
// Usage: rm [-r] [-f] <file>...
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	recursive = flag.Bool("r", false, "Remove directories recursively")
	force     = flag.Bool("f", false, "Ignore nonexistent files, never prompt")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: rm [-r] [-f] <file>...")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	exitCode := 0
	for _, path := range flag.Args() {
		if err := remove(path); err != nil {
			if !*force {
				fmt.Fprintln(os.Stderr, "rm:", err)
				exitCode = 1
			}
		}
	}
	os.Exit(exitCode)
}

func remove(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) && *force {
			return nil
		}
		return err
	}
	if info.IsDir() {
		if *recursive {
			return os.RemoveAll(path)
		}
		return fmt.Errorf("%s: is a directory (use -r)", path)
	}
	return os.Remove(path)
}
