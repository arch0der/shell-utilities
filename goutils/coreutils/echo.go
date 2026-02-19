package main

import (
	"fmt"
	"os"
	"strings"
)

func init() { register("echo", runEcho) }

func runEcho() {
	args := os.Args[1:]
	noNewline := false
	interpretEscapes := false
	i := 0
	for i < len(args) {
		a := args[i]
		if a == "-n" {
			noNewline = true
			i++
		} else if a == "-e" {
			interpretEscapes = true
			i++
		} else if a == "-E" {
			interpretEscapes = false
			i++
		} else {
			break
		}
	}
	out := strings.Join(args[i:], " ")
	if interpretEscapes {
		out = echoUnescape(out)
	}
	if noNewline {
		fmt.Print(out)
	} else {
		fmt.Println(out)
	}
}

