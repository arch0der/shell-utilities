// timer - stopwatch / elapsed time tracker with lap support
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func format(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

func main() {
	interactive := len(os.Args) == 1 || os.Args[1] == "-i"
	if interactive {
		start := time.Now()
		var laps []time.Duration
		fmt.Println("Timer started. Press ENTER for lap, 'q' to quit.")
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			elapsed := time.Since(start)
			switch {
			case line == "q" || line == "quit":
				fmt.Printf("Total: %s\n", format(elapsed)); return
			default:
				lapNum := len(laps) + 1
				lapTime := elapsed
				var split time.Duration
				if len(laps) > 0 { split = lapTime - laps[len(laps)-1] } else { split = lapTime }
				laps = append(laps, lapTime)
				fmt.Printf("Lap %-3d  Total: %s  Split: %s\n", lapNum, format(lapTime), format(split))
			}
		}
		return
	}

	// Non-interactive: time a command
	args := os.Args[1:]
	if args[0] == "-i" { args = args[1:] }
	if len(args) == 0 { fmt.Fprintln(os.Stderr, "usage: timer [command] [args...]"); os.Exit(1) }

	start := time.Now()
	// Just report; actual exec is left to shell wrapping
	fmt.Fprintf(os.Stderr, "Start: %s\n", start.Format("15:04:05.000"))
	fmt.Fprintf(os.Stderr, "Command: %s\n", strings.Join(args, " "))
}
