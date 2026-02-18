// xargs - Build and execute commands from stdin
// Usage: xargs [-n maxargs] [-I replace] [-P procs] <command> [args...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	maxArgs = flag.Int("n", 0, "Max arguments per command invocation (0=all)")
	replace = flag.String("I", "", "Replace string in command (e.g. -I{} cmd {} ...)")
	procs   = flag.Int("P", 1, "Max parallel processes")
	null    = flag.Bool("0", false, "Input items are separated by null, not whitespace")
	verbose = flag.Bool("t", false, "Print command before executing")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: xargs [-n N] [-I str] [-P N] [-0] [-t] <command> [args...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	baseCmd := flag.Args()
	if len(baseCmd) == 0 {
		baseCmd = []string{"echo"}
	}

	// Read input items
	var items []string
	scanner := bufio.NewScanner(os.Stdin)
	if *null {
		scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
			for i, b := range data {
				if b == 0 {
					return i + 1, data[:i], nil
				}
			}
			if atEOF && len(data) > 0 {
				return len(data), data, nil
			}
			return 0, nil, nil
		})
	}
	for scanner.Scan() {
		text := scanner.Text()
		if *null {
			items = append(items, text)
		} else {
			for _, field := range strings.Fields(text) {
				if field != "" {
					items = append(items, field)
				}
			}
		}
	}

	if len(items) == 0 {
		return
	}

	sem := make(chan struct{}, *procs)
	var wg sync.WaitGroup

	run := func(args []string) {
		defer wg.Done()
		sem <- struct{}{}
		defer func() { <-sem }()

		cmdArgs := make([]string, len(baseCmd))
		copy(cmdArgs, baseCmd)

		if *replace != "" {
			// Replace placeholder in command
			var out []string
			for _, a := range cmdArgs {
				out = append(out, strings.ReplaceAll(a, *replace, strings.Join(args, " ")))
			}
			cmdArgs = out
		} else {
			cmdArgs = append(cmdArgs, args...)
		}

		if *verbose {
			fmt.Fprintln(os.Stderr, strings.Join(cmdArgs, " "))
		}

		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "xargs:", err)
		}
	}

	if *replace != "" || *maxArgs == 1 {
		// One invocation per item
		for _, item := range items {
			wg.Add(1)
			go run([]string{item})
		}
	} else if *maxArgs > 0 {
		// Batch by maxArgs
		for i := 0; i < len(items); i += *maxArgs {
			end := i + *maxArgs
			if end > len(items) {
				end = len(items)
			}
			wg.Add(1)
			go run(items[i:end])
		}
	} else {
		// All at once
		wg.Add(1)
		go run(items)
	}

	wg.Wait()
}
