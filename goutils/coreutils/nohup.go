package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func init() { register("nohup", runNohup) }

func runNohup() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "nohup: missing operand")
		os.Exit(1)
	}

	// Redirect stdout to nohup.out if not a tty
	outFile := "nohup.out"
	if _, err := os.Stat(outFile); err != nil {
		f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err == nil {
			fmt.Fprintf(os.Stderr, "nohup: appending output to '%s'\n", outFile)
			_ = f
		}
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Ignore SIGHUP
	signal.Ignore(syscall.SIGHUP)

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "nohup: %v\n", err)
		os.Exit(1)
	}
}
