// dirname - Strip last component from file path
// Usage: dirname <path>...
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: dirname <path>...")
		os.Exit(1)
	}
	for _, path := range os.Args[1:] {
		fmt.Println(filepath.Dir(path))
	}
}
