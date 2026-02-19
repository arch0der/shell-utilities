// sync - Flush filesystem buffers to disk
// Usage: sync [file...]
package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	if len(os.Args) == 1 {
		// Sync everything
		syscall.Sync()
		return
	}

	// Sync specific files
	exitCode := 0
	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "sync:", err)
			exitCode = 1
			continue
		}
		if err := f.Sync(); err != nil {
			fmt.Fprintln(os.Stderr, "sync:", err)
			exitCode = 1
		}
		f.Close()
	}
	os.Exit(exitCode)
}
