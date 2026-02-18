// mv - Move or rename files
// Usage: mv <src> <dst>
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: mv <src> <dst>")
		os.Exit(1)
	}
	src := os.Args[1]
	dst := os.Args[2]

	// If dst is an existing directory, move src inside it
	if di, err := os.Stat(dst); err == nil && di.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	if err := os.Rename(src, dst); err != nil {
		fmt.Fprintln(os.Stderr, "mv:", err)
		os.Exit(1)
	}
}
