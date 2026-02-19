// pipe - run a sequence of commands as a pipeline, with error handling and logging
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: pipe [options] cmd1 '|' cmd2 '|' cmd3 ...
  -v    verbose: print each command before running
  -e    exit on first error (default: continue)
  Separate commands with literal '|' arguments.`)
	os.Exit(1)
}

func splitCmds(args []string) [][]string {
	var cmds [][]string
	cur := []string{}
	for _, a := range args {
		if a == "|" { cmds = append(cmds, cur); cur = nil } else { cur = append(cur, a) }
	}
	if len(cur) > 0 { cmds = append(cmds, cur) }
	return cmds
}

func main() {
	verbose := false
	exitOnErr := false
	args := os.Args[1:]
	filtered := args[:0]
	for _, a := range args {
		switch a {
		case "-v": verbose = true
		case "-e": exitOnErr = true
		default: filtered = append(filtered, a)
		}
	}

	cmds := splitCmds(filtered)
	if len(cmds) == 0 { usage() }

	// Build pipeline
	commands := make([]*exec.Cmd, len(cmds))
	for i, parts := range cmds {
		if len(parts) == 0 { usage() }
		if verbose { fmt.Fprintf(os.Stderr, "+ %s\n", strings.Join(parts, " ")) }
		commands[i] = exec.Command(parts[0], parts[1:]...)
	}

	// Chain stdin/stdout
	commands[0].Stdin = os.Stdin
	for i := 0; i < len(commands)-1; i++ {
		pipe, err := commands[i].StdoutPipe()
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		commands[i+1].Stdin = pipe
		commands[i].Stderr = os.Stderr
	}
	commands[len(commands)-1].Stdout = os.Stdout
	commands[len(commands)-1].Stderr = os.Stderr

	for _, c := range commands { c.Start() }
	code := 0
	for _, c := range commands {
		if err := c.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, "pipe:", err)
			if exitOnErr { os.Exit(1) }
			code = 1
		}
	}
	os.Exit(code)
}
