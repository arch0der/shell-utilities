package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func init() { register("timeout", runTimeout) }

func runTimeout() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "timeout: missing operand")
		os.Exit(1)
	}
	killAfter := time.Duration(0)
	signal := syscall.SIGTERM
	_ = signal
	files := []string{}
	i := 0
	for i < len(args) {
		a := args[i]
		if a == "-k" && i+1 < len(args) {
			i++
			d, err := parseDuration(args[i])
			if err == nil {
				killAfter = d
			}
		} else if a == "--preserve-status" || a == "--foreground" {
			// ignore
		} else {
			files = append(files, a)
		}
		i++
	}
	_ = killAfter
	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "timeout: missing command")
		os.Exit(1)
	}
	dur, err := parseDuration(files[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "timeout: invalid duration: %s\n", files[0])
		os.Exit(1)
	}
	cmd := exec.Command(files[1], files[2:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "timeout:", err)
		os.Exit(1)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
		}
	case <-time.After(dur):
		cmd.Process.Signal(syscall.SIGTERM)
		time.Sleep(100 * time.Millisecond)
		cmd.Process.Kill()
		os.Exit(124)
	}
}
