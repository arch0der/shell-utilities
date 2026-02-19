package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func main() {
	args := os.Args[1:]
	all := false
	save := false

	for _, a := range args {
		if a == "-a" || a == "--all" {
			all = true
		} else if a == "-g" || a == "--save" {
			save = true
		}
	}

	var termios syscall.Termios
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(0),
		syscall.TCGETS, uintptr(unsafe.Pointer(&termios))); err != 0 {
		fmt.Fprintln(os.Stderr, "stty: not a terminal")
		os.Exit(1)
	}

	if save {
		fmt.Printf("%x:%x:%x:%x", termios.Iflag, termios.Oflag, termios.Cflag, termios.Lflag)
		fmt.Println()
		return
	}

	if all || len(args) == 0 {
		speed := "38400"
		fmt.Printf("speed %s baud; rows %d; columns %d;\n", speed, 24, 80)
		fmt.Printf("intr = ^C; quit = ^\\; erase = ^?; kill = ^U;\n")
		if all {
			fmt.Printf("iflags: icrnl ixon\n")
			fmt.Printf("oflags: opost onlcr\n")
			fmt.Printf("lflags: isig icanon echo echoe echok\n")
		}
		return
	}

	// Handle settings
	for _, a := range args {
		switch strings.TrimPrefix(a, "-") {
		case "echo", "-echo", "icanon", "-icanon":
			// Would set termios flags here
		}
	}
}
