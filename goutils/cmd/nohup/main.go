// nohup - Run a command immune to hangups, with output to nohup.out
// Usage: nohup <command> [args...]
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: nohup <command> [args...]")
		os.Exit(1)
	}

	// Ignore SIGHUP
	signal.Ignore(syscall.SIGHUP)

	// Redirect stdout to nohup.out if it's a terminal
	outFile := os.Stdout
	if isTerminal(os.Stdout.Fd()) {
		f, err := os.OpenFile("nohup.out", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			// Try home directory
			home, _ := os.UserHomeDir()
			f, err = os.OpenFile(home+"/nohup.out", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, "nohup: cannot open output file")
				os.Exit(1)
			}
		}
		fmt.Fprintf(os.Stderr, "nohup: appending output to '%s'\n", f.Name())
		outFile = f
		defer f.Close()
	}

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "nohup:", err)
		os.Exit(1)
	}
}

func isTerminal(fd uintptr) bool {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TIOCGWINSZ, 0)
	return err == 0
}
