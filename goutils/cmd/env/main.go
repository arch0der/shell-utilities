// env - Print or set environment variables
// Usage: env [NAME=VALUE ...] [command [args...]]
package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func main() {
	args := os.Args[1:]

	// Separate NAME=VALUE pairs from command
	var pairs []string
	cmdStart := len(args)
	for i, a := range args {
		if !strings.Contains(a, "=") {
			cmdStart = i
			break
		}
		pairs = append(pairs, a)
	}

	// Set env vars
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		os.Setenv(parts[0], parts[1])
	}

	if cmdStart < len(args) {
		// Run command with modified environment
		cmd := exec.Command(args[cmdStart], args[cmdStart+1:]...)
		cmd.Env = os.Environ()
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "env:", err)
			os.Exit(1)
		}
		return
	}

	// Print environment
	env := os.Environ()
	sort.Strings(env)
	for _, e := range env {
		fmt.Println(e)
	}
}
