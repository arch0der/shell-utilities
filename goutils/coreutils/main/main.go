// coreutils - multi-call binary implementing GNU coreutils in Go
// Usage: coreutils <command> [args...]  OR  symlink/copy binary as <command>
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

var registry = map[string]func(){}

func main() {
	name := filepath.Base(os.Args[0])
	args := os.Args[1:]

	if name == "coreutils" || name == "coreutils.exe" {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: coreutils <command> [args...]")
			fmt.Fprintln(os.Stderr, "\nAvailable commands:")
			names := make([]string, 0, len(registry))
			for k := range registry {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, k := range names {
				fmt.Fprintf(os.Stderr, "  %s\n", k)
			}
			os.Exit(1)
		}
		name = args[0]
		args = args[1:]
	}
	os.Args = append([]string{name}, args...)

	fn, ok := registry[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "coreutils: %s: command not found\n", name)
		os.Exit(1)
	}
	fn()
}
