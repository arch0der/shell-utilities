package main

import (
	"fmt"
	"os"
)

func init() { register("chcon", runChcon) }

func runChcon() {
	// SELinux context change - stub (requires SELinux)
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "chcon: missing operand")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "chcon: SELinux not supported on this platform")
	os.Exit(1)
}
