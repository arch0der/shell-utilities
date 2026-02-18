// touch - Create empty files or update timestamps
// Usage: touch <file>...
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: touch <file>...")
		os.Exit(1)
	}
	now := time.Now()
	exitCode := 0
	for _, path := range os.Args[1:] {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "touch:", err)
			exitCode = 1
			continue
		}
		f.Close()
		if err := os.Chtimes(path, now, now); err != nil {
			fmt.Fprintln(os.Stderr, "touch:", err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
