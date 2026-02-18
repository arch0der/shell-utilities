// which - Locate commands in PATH
// Usage: which [-a] <command>...
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var all = flag.Bool("a", false, "Print all matching paths")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: which [-a] <command>...")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	pathEnv := os.Getenv("PATH")
	dirs := strings.Split(pathEnv, string(os.PathListSeparator))

	exitCode := 0
	for _, cmd := range flag.Args() {
		found := false
		for _, dir := range dirs {
			full := filepath.Join(dir, cmd)
			info, err := os.Stat(full)
			if err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
				fmt.Println(full)
				found = true
				if !*all {
					break
				}
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "%s not found\n", cmd)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
