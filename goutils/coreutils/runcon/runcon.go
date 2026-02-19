package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "runcon: SELinux not supported on this platform")
	os.Exit(1)
}
