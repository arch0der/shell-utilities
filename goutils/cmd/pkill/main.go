// pkill - Kill processes matching a pattern
// Usage: pkill [-signal] [-x] [-u user] <pattern>
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var (
	sigName = flag.String("s", "TERM", "Signal to send")
	exact   = flag.Bool("x", false, "Exact match")
	user    = flag.String("u", "", "Match by UID")
)

var signals = map[string]syscall.Signal{
	"HUP": syscall.SIGHUP, "INT": syscall.SIGINT, "QUIT": syscall.SIGQUIT,
	"KILL": syscall.SIGKILL, "TERM": syscall.SIGTERM, "USR1": syscall.SIGUSR1,
	"USR2": syscall.SIGUSR2, "STOP": syscall.SIGSTOP, "CONT": syscall.SIGCONT,
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: pkill [-s signal] [-x] [-u user] <pattern>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	sig, ok := signals[strings.ToUpper(*sigName)]
	if !ok {
		fmt.Fprintf(os.Stderr, "pkill: unknown signal %s\n", *sigName)
		os.Exit(1)
	}

	pattern := flag.Arg(0)
	if *exact {
		pattern = "^" + pattern + "$"
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pkill: invalid pattern:", err)
		os.Exit(1)
	}

	entries, _ := os.ReadDir("/proc")
	killed := 0
	myPid := os.Getpid()

	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil || pid == myPid {
			continue
		}
		base := filepath.Join("/proc", e.Name())
		status := readFileStr(base + "/status")
		name := extractField(status, "Name")
		uid := extractUID(status)

		if *user != "" && uid != *user {
			continue
		}
		if !re.MatchString(name) {
			continue
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			continue
		}
		if err := proc.Signal(sig); err == nil {
			killed++
		}
	}

	if killed == 0 {
		os.Exit(1)
	}
}

func readFileStr(path string) string {
	b, _ := os.ReadFile(path)
	return string(b)
}

func extractField(status, field string) string {
	for _, line := range strings.Split(status, "\n") {
		if strings.HasPrefix(line, field+":") {
			return strings.TrimSpace(strings.TrimPrefix(line, field+":"))
		}
	}
	return ""
}

func extractUID(status string) string {
	for _, line := range strings.Split(status, "\n") {
		if strings.HasPrefix(line, "Uid:") {
			f := strings.Fields(line)
			if len(f) >= 2 {
				return f[1]
			}
		}
	}
	return ""
}
