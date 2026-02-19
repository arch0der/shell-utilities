// ratelimit - rate-limit stdin lines: pass through N lines per second (or per interval)
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: ratelimit [options]
  -r <rate>     lines per second (default: 1)
  -n <n>        lines per interval (default: 1)
  -i <interval> interval duration e.g. 500ms, 1s, 2m (default: 1s)
  -v            verbose: print timing info to stderr`)
	os.Exit(1)
}

func parseDuration(s string) (time.Duration, error) { return time.ParseDuration(s) }

func main() {
	linesPerInterval := 1
	interval := time.Second
	verbose := false

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-r":
			i++; r, err := strconv.ParseFloat(args[i], 64)
			if err != nil { usage() }
			interval = time.Duration(float64(time.Second) / r)
		case "-n": i++; linesPerInterval, _ = strconv.Atoi(args[i])
		case "-i":
			i++; d, err := parseDuration(args[i])
			if err != nil { fmt.Fprintln(os.Stderr, "ratelimit: bad interval:", err); os.Exit(1) }
			interval = d
		case "-v": verbose = true
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
		}
	}

	sc := bufio.NewScanner(os.Stdin)
	batch := 0
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	buf := []string{}
	done := make(chan struct{})

	go func() {
		for sc.Scan() { buf = append(buf, sc.Text()) }
		close(done)
	}()

	idx := 0
	for {
		select {
		case t := <-ticker.C:
			for i := 0; i < linesPerInterval && idx < len(buf); i++ {
				fmt.Println(buf[idx]); idx++; batch++
			}
			if verbose { fmt.Fprintf(os.Stderr, "[%s] sent %d lines total\n", t.Format("15:04:05.000"), batch) }
		case <-done:
			for idx < len(buf) { fmt.Println(buf[idx]); idx++ }
			return
		}
		if idx >= len(buf) { select { case <-done: return; default: } }
	}
}
