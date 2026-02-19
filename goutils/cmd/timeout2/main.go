// timeout2 - alias for timeout with different name to avoid conflict
// Usage: timeout2 <duration> <cmd> [args...]
package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: timeout2 <duration> <cmd> [args...]")
		os.Exit(1)
	}
	dur, err := time.ParseDuration(os.Args[1])
	if err != nil {
		f := 0.0
		fmt.Sscanf(os.Args[1], "%f", &f)
		dur = time.Duration(f * float64(time.Second))
	}
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			if e, ok := err.(*exec.ExitError); ok {
				os.Exit(e.ExitCode())
			}
			os.Exit(1)
		}
	case <-time.After(dur):
		cmd.Process.Signal(syscall.SIGTERM)
		time.Sleep(time.Second)
		cmd.Process.Kill()
		fmt.Fprintf(os.Stderr, "timeout: timed out after %s\n", dur)
		os.Exit(124)
	}
}
