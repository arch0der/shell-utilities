package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	args := os.Args[1:]
	cmdStart := 0
	for i, a := range args {
		if !strings.HasPrefix(a, "-") {
			cmdStart = i
			break
		}
		// -i, -o, -e flags with values
		if (a == "-i" || a == "-o" || a == "-e") && i+1 < len(args) {
			i++
		}
	}
	if cmdStart >= len(args) {
		fmt.Fprintln(os.Stderr, "stdbuf: missing command")
		os.Exit(1)
	}
	cmd := exec.Command(args[cmdStart], args[cmdStart+1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "stdbuf:", err)
		os.Exit(1)
	}
}
