// fwatch - Watch files for changes and run a command.
//
// Usage:
//
//	fwatch [OPTIONS] FILE [FILE...] -- COMMAND [ARGS...]
//	fwatch [OPTIONS] PATTERN -- COMMAND [ARGS...]
//
// Options:
//
//	-i DUR    Polling interval (default: 500ms)
//	-d        Debounce: wait for changes to settle before running (default: 200ms)
//	-1        Run once then exit (don't keep watching)
//	-c        Clear screen before each run
//	-r        Recursive watch (directories)
//	-p PAT    Pattern filter (glob, e.g. "*.go")
//	-q        Quiet: don't print change notifications
//	-n        Dry run: print what would run but don't execute
//
// Examples:
//
//	fwatch main.go -- go run main.go
//	fwatch -p "*.go" . -- go test ./...
//	fwatch -c -r src/ -- npm run build
//	fwatch -i 1s config.yaml -- systemctl reload app
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	interval  = flag.Duration("i", 500*time.Millisecond, "poll interval")
	debounce  = flag.Duration("d", 200*time.Millisecond, "debounce wait")
	once      = flag.Bool("1", false, "run once")
	clear     = flag.Bool("c", false, "clear screen")
	recursive = flag.Bool("r", false, "recursive")
	pattern   = flag.String("p", "", "file pattern")
	quiet     = flag.Bool("q", false, "quiet")
	dryRun    = flag.Bool("n", false, "dry run")
)

type FileState struct {
	mtime time.Time
	size  int64
}

func getState(paths []string) map[string]FileState {
	state := make(map[string]FileState)
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if info.IsDir() && *recursive {
			filepath.Walk(p, func(wp string, wi os.FileInfo, err error) error {
				if err != nil || wi.IsDir() {
					return nil
				}
				if *pattern != "" {
					matched, _ := filepath.Match(*pattern, filepath.Base(wp))
					if !matched {
						return nil
					}
				}
				state[wp] = FileState{wi.ModTime(), wi.Size()}
				return nil
			})
		} else if !info.IsDir() {
			state[p] = FileState{info.ModTime(), info.Size()}
		}
	}
	return state
}

func stateChanged(a, b map[string]FileState) (bool, string) {
	for k, v := range b {
		if prev, ok := a[k]; !ok {
			return true, k + " (new)"
		} else if v.mtime != prev.mtime || v.size != prev.size {
			return true, k + " (modified)"
		}
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			return true, k + " (deleted)"
		}
	}
	return false, ""
}

func runCmd(cmd []string) {
	if *dryRun {
		fmt.Printf("[dry-run] %s\n", strings.Join(cmd, " "))
		return
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Run()
}

func main() {
	flag.Parse()
	args := flag.Args()

	// Split on "--"
	var watchPaths, command []string
	for i, a := range args {
		if a == "--" {
			watchPaths = args[:i]
			command = args[i+1:]
			break
		}
	}
	if len(command) == 0 {
		// No -- separator: last arg that looks like a command
		fmt.Fprintln(os.Stderr, "usage: fwatch [OPTIONS] FILE... -- COMMAND [ARGS...]")
		os.Exit(1)
	}
	if len(watchPaths) == 0 {
		watchPaths = []string{"."}
	}

	state := getState(watchPaths)

	// Run once immediately
	if *clear {
		fmt.Print("\033[2J\033[H")
	}
	runCmd(command)

	if *once {
		return
	}

	var pendingChange string
	var lastChange time.Time

	fmt.Fprintf(os.Stderr, "watching %s...\n", strings.Join(watchPaths, ", "))

	for {
		time.Sleep(*interval)

		newState := getState(watchPaths)
		changed, what := stateChanged(state, newState)

		if changed {
			state = newState
			pendingChange = what
			lastChange = time.Now()
		}

		if pendingChange != "" && time.Since(lastChange) >= *debounce {
			if !*quiet {
				fmt.Fprintf(os.Stderr, "\n[fwatch] changed: %s\n", pendingChange)
			}
			pendingChange = ""
			if *clear {
				fmt.Print("\033[2J\033[H")
			}
			runCmd(command)
		}
	}
}
