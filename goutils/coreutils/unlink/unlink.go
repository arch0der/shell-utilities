package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "unlink: missing operand")
		os.Exit(1)
	}
	if err := os.Remove(args[0]); err != nil {
		fmt.Fprintln(os.Stderr, "unlink:", err)
		os.Exit(1)
	}
}
