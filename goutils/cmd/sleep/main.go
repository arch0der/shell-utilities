// sleep - Delay for a specified amount of time
// Usage: sleep <duration>...
// Supports: 1s, 1m, 1h, 1.5, 500ms, or plain seconds
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: sleep <duration>...")
		fmt.Fprintln(os.Stderr, "Examples: sleep 1  sleep 1.5  sleep 2s  sleep 1m30s")
		os.Exit(1)
	}

	var total time.Duration
	for _, arg := range os.Args[1:] {
		d, err := parseDuration(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sleep: invalid time interval '%s'\n", arg)
			os.Exit(1)
		}
		total += d
	}

	time.Sleep(total)
}

func parseDuration(s string) (time.Duration, error) {
	// Try Go's native format first (1s, 1m, 500ms, etc.)
	d, err := time.ParseDuration(s)
	if err == nil {
		return d, nil
	}

	// Try plain float (seconds)
	f, err2 := strconv.ParseFloat(s, 64)
	if err2 != nil {
		return 0, err
	}
	return time.Duration(f * float64(time.Second)), nil
}
