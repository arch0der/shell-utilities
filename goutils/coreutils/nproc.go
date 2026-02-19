package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

func init() { register("nproc", runNproc) }

func runNproc() {
	args := os.Args[1:]
	all := false
	ignore := 0
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--all" {
			all = true
		} else if a == "--ignore" && i+1 < len(args) {
			i++
			ignore, _ = strconv.Atoi(args[i])
		}
	}
	_ = all
	n := runtime.NumCPU() - ignore
	if n < 1 {
		n = 1
	}
	fmt.Println(n)
}
