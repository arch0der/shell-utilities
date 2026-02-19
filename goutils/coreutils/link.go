package main

import (
	"fmt"
	"os"
)

func init() { register("link", runLink) }

func runLink() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "link: missing operand")
		os.Exit(1)
	}
	if err := os.Link(args[0], args[1]); err != nil {
		fmt.Fprintln(os.Stderr, "link:", err)
		os.Exit(1)
	}
}
