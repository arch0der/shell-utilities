// ps - List running processes
// Usage: ps [-a]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var all = flag.Bool("a", false, "Show all processes")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: ps [-a]")
		flag.PrintDefaults()
	}
	flag.Parse()

	myPid := os.Getpid()
	myUid := os.Getuid()

	entries, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ps:", err)
		os.Exit(1)
	}

	fmt.Printf("%-8s %-8s %-8s %s\n", "PID", "UID", "STAT", "CMD")

	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		if pid == myPid {
			continue
		}

		// Read owner
		info, err := os.Stat(filepath.Join("/proc", e.Name()))
		if err != nil {
			continue
		}
		_ = info

		// Read status
		statusBytes, err := os.ReadFile(filepath.Join("/proc", e.Name(), "status"))
		if err != nil {
			continue
		}
		status := parseStatus(string(statusBytes))

		uid, _ := strconv.Atoi(status["Uid"])
		if !*all && uid != myUid {
			continue
		}

		// Read cmdline
		cmdBytes, _ := os.ReadFile(filepath.Join("/proc", e.Name(), "cmdline"))
		cmd := strings.ReplaceAll(string(cmdBytes), "\x00", " ")
		if cmd == "" {
			cmd = "[" + status["Name"] + "]"
		}
		cmd = strings.TrimSpace(cmd)

		fmt.Printf("%-8d %-8d %-8s %s\n", pid, uid, status["State"], cmd)
	}
}

func parseStatus(s string) map[string]string {
	m := map[string]string{}
	for _, line := range strings.Split(s, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			// Uid field has multiple values; take first
			if key == "Uid" {
				val = strings.Fields(val)[0]
			}
			m[key] = val
		}
	}
	return m
}
