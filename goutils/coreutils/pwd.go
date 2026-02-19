package main

import (
	"fmt"
	"os"
)

func init() { register("pwd", runPwd) }

func runPwd() {
	args := os.Args[1:]
	logical := true
	for _, a := range args {
		if a == "-P" || a == "--physical" {
			logical = false
		} else if a == "-L" || a == "--logical" {
			logical = true
		}
	}
	_ = logical
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pwd:", err)
		os.Exit(1)
	}
	fmt.Println(dir)
}
