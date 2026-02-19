package main

import (
	"fmt"
	"os"
)

func init() { register("hostname", runHostname) }

func runHostname() {
	args := os.Args[1:]
	if len(args) > 0 && !func() bool {
		for _, a := range args {
			if a[0] == '-' {
				return true
			}
		}
		return false
	}() {
		// Set hostname
		if err := setHostname(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "hostname: %v\n", err)
			os.Exit(1)
		}
		return
	}
	h, err := os.Hostname()
	if err != nil {
		fmt.Fprintln(os.Stderr, "hostname:", err)
		os.Exit(1)
	}
	fmt.Println(h)
}
