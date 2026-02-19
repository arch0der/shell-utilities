package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	args := os.Args[1:]
	ignore := false
	unsets := []string{}
	sets := []string{}
	cmdStart := len(args)
	nullSep := false

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-i" || a == "--ignore-environment":
			ignore = true
		case a == "-0" || a == "--null":
			nullSep = true
		case a == "-u" && i+1 < len(args):
			i++
			unsets = append(unsets, args[i])
		case strings.HasPrefix(a, "-u"):
			unsets = append(unsets, a[2:])
		case strings.Contains(a, "="):
			sets = append(sets, a)
		default:
			cmdStart = i
			i = len(args) // break
		}
	}

	env := os.Environ()
	if ignore {
		env = []string{}
	}
	for _, u := range unsets {
		newEnv := env[:0]
		for _, e := range env {
			if !strings.HasPrefix(e, u+"=") {
				newEnv = append(newEnv, e)
			}
		}
		env = newEnv
	}
	env = append(env, sets...)

	if cmdStart >= len(args) {
		sep := "\n"
		if nullSep {
			sep = "\x00"
		}
		for _, e := range env {
			fmt.Print(e + sep)
		}
		return
	}

	cmd := exec.Command(args[cmdStart], args[cmdStart+1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "env: %v\n", err)
		os.Exit(126)
	}
}
