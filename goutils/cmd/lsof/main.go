// lsof - List open files (Linux /proc based)
// Usage: lsof [-p pid] [-u user] [-c name]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	pid      = flag.Int("p", 0, "Show only files for this PID")
	userName = flag.String("u", "", "Show only files for this UID")
	cmdName  = flag.String("c", "", "Show only files for commands matching name")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: lsof [-p pid] [-u uid] [-c name]")
		flag.PrintDefaults()
	}
	flag.Parse()

	entries, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Fprintln(os.Stderr, "lsof:", err)
		os.Exit(1)
	}

	fmt.Printf("%-8s %-8s %-20s %s\n", "PID", "TYPE", "CMD", "NAME")

	for _, e := range entries {
		p, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		if *pid != 0 && p != *pid {
			continue
		}

		base := filepath.Join("/proc", e.Name())
		cmd := readFile(filepath.Join(base, "comm"))

		if *cmdName != "" && !strings.Contains(cmd, *cmdName) {
			continue
		}

		// Check uid
		if *userName != "" {
			status := readFile(filepath.Join(base, "status"))
			uid := extractUID(status)
			if uid != *userName {
				continue
			}
		}

		// Read fd directory
		fdDir := filepath.Join(base, "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}
		for _, fd := range fds {
			target, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			fdType := "REG"
			if strings.HasPrefix(target, "pipe:") {
				fdType = "PIPE"
			} else if strings.HasPrefix(target, "socket:") {
				fdType = "SOCK"
			} else if strings.HasPrefix(target, "/dev/") {
				fdType = "CHR"
			}
			fmt.Printf("%-8d %-8s %-20s %s\n", p, fdType, strings.TrimSpace(cmd), target)
		}
	}
}

func readFile(path string) string {
	b, _ := os.ReadFile(path)
	return strings.TrimSpace(string(b))
}

func extractUID(status string) string {
	for _, line := range strings.Split(status, "\n") {
		if strings.HasPrefix(line, "Uid:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[1]
			}
		}
	}
	return ""
}
