package main

import (
	"fmt"
	"os"
	"strings"
)

func init() { register("yes", runYes) }

func runYes() {
	args := os.Args[1:]
	msg := "y"
	if len(args) > 0 {
		msg = strings.Join(args, " ")
	}
	for {
		fmt.Fprintln(os.Stdout, msg)
		if err := os.Stdout.Sync(); err != nil {
			os.Exit(0)
		}
	}
}
