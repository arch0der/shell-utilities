// timeout - Run a command with a time limit
// Usage: timeout [-s signal] <duration> <command> [args...]
// Duration format: 10s, 1m, 1h, 500ms
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var sigName = flag.String("s", "TERM", "Signal to send on timeout (TERM, KILL, INT, etc.)")

var signals = map[string]syscall.Signal{
	"HUP":  syscall.SIGHUP,
	"INT":  syscall.SIGINT,
	"QUIT": syscall.SIGQUIT,
	"KILL": syscall.SIGKILL,
	"TERM": syscall.SIGTERM,
	"USR1": syscall.SIGUSR1,
	"USR2": syscall.SIGUSR2,
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: timeout [-s signal] <duration> <command> [args...]")
		fmt.Fprintln(os.Stderr, "Duration examples: 10s, 1m30s, 500ms")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	dur, err := time.ParseDuration(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "timeout: invalid duration %q: %v\n", flag.Arg(0), err)
		os.Exit(1)
	}

	sig, ok := signals[*sigName]
	if !ok {
		fmt.Fprintf(os.Stderr, "timeout: unknown signal %s\n", *sigName)
		os.Exit(1)
	}

	cmd := exec.Command(flag.Arg(1), flag.Args()[2:]...)
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
			os.Exit(1)
		}
	case <-time.After(dur):
		cmd.Process.Signal(sig)
		// If KILL was not the signal, give a grace period then force kill
		if sig != syscall.SIGKILL {
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				cmd.Process.Kill()
			}
		}
		fmt.Fprintf(os.Stderr, "timeout: command timed out after %s\n", dur)
		os.Exit(124) // Standard timeout exit code
	}
}
