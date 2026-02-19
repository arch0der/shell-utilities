package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func init() { register("chroot", runChroot) }

func runChroot() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "chroot: missing operand")
		os.Exit(1)
	}
	newRoot := args[0]
	cmdArgs := args[1:]
	if len(cmdArgs) == 0 {
		cmdArgs = []string{"/bin/sh", "-i"}
	}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: newRoot,
	}
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "chroot: %v\n", err)
		os.Exit(1)
	}
}
