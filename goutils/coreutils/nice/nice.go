package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	args := os.Args[1:]
	adjustment := 10
	cmdArgs := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-n" && i+1 < len(args) {
			i++
			adjustment, _ = strconv.Atoi(args[i])
		} else if strings.HasPrefix(a, "-n") {
			adjustment, _ = strconv.Atoi(a[2:])
		} else if len(a) > 1 && a[0] == '-' {
			n, err := strconv.Atoi(a[1:])
			if err == nil {
				adjustment = -n
			} else {
				cmdArgs = append(cmdArgs, a)
			}
		} else {
			cmdArgs = append(cmdArgs, args[i:]...)
			break
		}
	}

	if len(cmdArgs) == 0 {
		// Print current niceness
		nice, _ := syscall.Getpriority(0, 0)
		fmt.Println(nice)
		return
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	// Set niceness before exec
	syscall.Setpriority(0, 0, adjustment)

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "nice: %v\n", err)
		os.Exit(1)
	}
}
