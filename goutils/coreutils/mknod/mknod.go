package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	args := os.Args[1:]
	mode := uint32(0666)
	files := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-m" && i+1 < len(args) {
			i++
			m, _ := strconv.ParseUint(args[i], 8, 32)
			mode = uint32(m)
		} else if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}
	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "mknod: missing operand")
		os.Exit(1)
	}
	name := files[0]
	devType := files[1]
	major, minor := uint64(0), uint64(0)
	if len(files) >= 4 {
		major, _ = strconv.ParseUint(files[2], 10, 32)
		minor, _ = strconv.ParseUint(files[3], 10, 32)
	}
	var nodeType uint32
	switch devType {
	case "b":
		nodeType = syscall.S_IFBLK
	case "c", "u":
		nodeType = syscall.S_IFCHR
	case "p":
		if err := syscall.Mkfifo(name, mode); err != nil {
			fmt.Fprintf(os.Stderr, "mknod: %v\n", err)
			os.Exit(1)
		}
		return
	default:
		fmt.Fprintf(os.Stderr, "mknod: invalid device type '%s'\n", devType)
		os.Exit(1)
	}
	dev := major*256 + minor
	if err := syscall.Mknod(name, nodeType|mode, int(dev)); err != nil {
		fmt.Fprintf(os.Stderr, "mknod: %v\n", err)
		os.Exit(1)
	}
}
