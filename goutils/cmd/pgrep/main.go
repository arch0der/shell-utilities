// pgrep - Find processes by name or attribute
// Usage: pgrep [-l] [-x] [-u user] <pattern>
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	list    = flag.Bool("l", false, "Also print process name")
	exact   = flag.Bool("x", false, "Exact match (whole name)")
	user    = flag.String("u", "", "Match by UID")
	newest  = flag.Bool("n", false, "Select only newest")
	oldest  = flag.Bool("o", false, "Select only oldest")
	count   = flag.Bool("c", false, "Count matching processes")
)

type Process struct {
	pid  int
	name string
	uid  string
	start int64
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: pgrep [-l] [-x] [-u user] [-n] [-o] [-c] <pattern>")
		flag.PrintDefaults()
	}
	flag.Parse()

	pattern := ""
	if flag.NArg() > 0 {
		pattern = flag.Arg(0)
	}

	entries, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Fprintln(os.Stderr, "pgrep:", err)
		os.Exit(1)
	}

	myPid := os.Getpid()
	var re *regexp.Regexp
	if pattern != "" {
		if *exact {
			pattern = "^" + pattern + "$"
		}
		re, err = regexp.Compile(pattern)
		if err != nil {
			fmt.Fprintln(os.Stderr, "pgrep: invalid pattern:", err)
			os.Exit(1)
		}
	}

	var procs []Process
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil || pid == myPid {
			continue
		}
		base := filepath.Join("/proc", e.Name())
		status := readFileStr(filepath.Join(base, "status"))
		name := extractField(status, "Name")
		uid := extractUID(status)

		if *user != "" && uid != *user {
			continue
		}
		if re != nil && !re.MatchString(name) {
			continue
		}

		// Get start time from stat
		stat := readFileStr(filepath.Join(base, "stat"))
		startTime := extractStartTime(stat)

		procs = append(procs, Process{pid, name, uid, startTime})
	}

	if len(procs) == 0 {
		os.Exit(1)
	}

	if *count {
		fmt.Println(len(procs))
		return
	}

	if *newest && len(procs) > 1 {
		best := procs[0]
		for _, p := range procs[1:] {
			if p.start > best.start {
				best = p
			}
		}
		procs = []Process{best}
	}
	if *oldest && len(procs) > 1 {
		best := procs[0]
		for _, p := range procs[1:] {
			if p.start < best.start {
				best = p
			}
		}
		procs = []Process{best}
	}

	for _, p := range procs {
		if *list {
			fmt.Printf("%d %s\n", p.pid, p.name)
		} else {
			fmt.Println(p.pid)
		}
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

func extractStartTime(stat string) int64 {
	fields := strings.Fields(stat)
	if len(fields) >= 22 {
		n, _ := strconv.ParseInt(fields[21], 10, 64)
		return n
	}
	return 0
}
