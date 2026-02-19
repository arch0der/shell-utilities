package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func init() { register("mkfifo", runMkfifo) }

func runMkfifo() {
	args := os.Args[1:]
	mode := uint32(0666)
	files := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-m" && i+1 < len(args) {
			i++
			m, _ := strconv.ParseUint(args[i], 8, 32)
			mode = uint32(m)
		} else if strings.HasPrefix(a, "-m") {
			m, _ := strconv.ParseUint(a[2:], 8, 32)
			mode = uint32(m)
		} else if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}
	exitCode := 0
	for _, f := range files {
		if err := syscall.Mkfifo(f, mode); err != nil {
			fmt.Fprintf(os.Stderr, "mkfifo: %s: %v\n", f, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
