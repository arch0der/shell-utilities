package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "sleep: missing operand")
		os.Exit(1)
	}
	var total time.Duration
	for _, a := range args {
		d, err := parseDuration(a)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sleep: invalid time interval '%s'\n", a)
			os.Exit(1)
		}
		total += d
	}
	time.Sleep(total)
}
