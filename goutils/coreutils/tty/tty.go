package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func main() {
	args := os.Args[1:]
	silent := false
	for _, a := range args {
		if a == "-s" || a == "--silent" || a == "--quiet" {
			silent = true
		}
	}
	// Check if stdin is a tty by getting the tty name
	name := ttyName()
	if name == "" {
		if !silent {
			fmt.Println("not a tty")
		}
		os.Exit(1)
	}
	if !silent {
		fmt.Println(name)
	}
}

func ttyName() string {
	// Try /proc/self/fd/0 -> readlink
	link, err := os.Readlink("/proc/self/fd/0")
	if err == nil {
		return link
	}
	// Fallback: check if it's a terminal via ioctl
	var t [1]uint32
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, 0, syscall.TIOCGPGRP, uintptr(unsafe.Pointer(&t[0]))); errno == 0 {
		return "/dev/tty"
	}
	return ""
}
