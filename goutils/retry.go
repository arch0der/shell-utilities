// retry - Retry a command on failure with backoff.
//
// Usage:
//
//	retry [OPTIONS] -- COMMAND [ARGS...]
//
// Options:
//
//	-n N       Max attempts (default: 3)
//	-w DUR     Initial wait between retries (default: 1s)
//	-b FACTOR  Backoff factor (default: 2.0)
//	-max DUR   Maximum wait between retries (default: 30s)
//	-j DUR     Jitter: random +/- duration added to wait
//	-q         Quiet: suppress retry messages
//	-e CODES   Retry only on these exit codes (comma-separated, default: any non-zero)
//	-t DUR     Total timeout for all attempts
//
// Examples:
//
//	retry -- curl https://api.example.com/data
//	retry -n 5 -w 2s -b 2 -- wget https://example.com/file.zip
//	retry -n 10 -w 100ms -max 5s -- ./check-ready.sh
//	retry -t 60s -- ssh host "command"
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	maxAttempts = flag.Int("n", 3, "max attempts")
	wait        = flag.Duration("w", time.Second, "initial wait")
	backoff     = flag.Float64("b", 2.0, "backoff factor")
	maxWait     = flag.Duration("max", 30*time.Second, "max wait")
	jitter      = flag.Duration("j", 0, "jitter")
	quiet       = flag.Bool("q", false, "quiet")
	exitCodes   = flag.String("e", "", "retry on these exit codes")
	totalTime   = flag.Duration("t", 0, "total timeout")
)

func parseExitCodes(s string) map[int]bool {
	if s == "" {
		return nil
	}
	m := make(map[int]bool)
	for _, p := range strings.Split(s, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err == nil {
			m[n] = true
		}
	}
	return m
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: retry [OPTIONS] -- COMMAND [ARGS...]")
		os.Exit(1)
	}

	allowed := parseExitCodes(*exitCodes)
	var deadline time.Time
	if *totalTime > 0 {
		deadline = time.Now().Add(*totalTime)
	}

	currentWait := *wait
	var lastCode int

	for attempt := 1; attempt <= *maxAttempts; attempt++ {
		if !deadline.IsZero() && time.Now().After(deadline) {
			fmt.Fprintf(os.Stderr, "retry: total timeout exceeded\n")
			os.Exit(lastCode)
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err == nil {
			os.Exit(0)
		}

		code := 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		}
		lastCode = code

		// Check if we should retry this code
		if allowed != nil && !allowed[code] {
			if !*quiet {
				fmt.Fprintf(os.Stderr, "retry: exit code %d not in retry list, giving up\n", code)
			}
			os.Exit(code)
		}

		if attempt >= *maxAttempts {
			break
		}

		sleepDur := currentWait
		if *jitter > 0 {
			j := time.Duration(rand.Int63n(int64(*jitter)*2)) - *jitter
			sleepDur += j
			if sleepDur < 0 {
				sleepDur = 0
			}
		}

		if !*quiet {
			fmt.Fprintf(os.Stderr, "retry: attempt %d/%d failed (exit %d), waiting %s...\n",
				attempt, *maxAttempts, code, sleepDur.Round(time.Millisecond))
		}

		// Check deadline before sleeping
		if !deadline.IsZero() {
			remaining := time.Until(deadline)
			if sleepDur > remaining {
				sleepDur = remaining
			}
		}

		time.Sleep(sleepDur)

		// Apply backoff
		next := float64(currentWait) * *backoff
		currentWait = time.Duration(math.Min(next, float64(*maxWait)))
	}

	if !*quiet {
		fmt.Fprintf(os.Stderr, "retry: all %d attempts failed\n", *maxAttempts)
	}
	os.Exit(lastCode)
}
