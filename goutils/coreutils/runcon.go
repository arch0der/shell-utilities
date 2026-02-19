package main

import (
	"fmt"
	"os"
)

func init() { register("runcon", runRuncon) }

func runRuncon() {
	fmt.Fprintln(os.Stderr, "runcon: SELinux not supported on this platform")
	os.Exit(1)
}
