package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

func init() { register("uptime", runUptime) }

func runUptime() {
	args := os.Args[1:]
	pretty := false
	since := false
	for _, a := range args {
		if a == "-p" || a == "--pretty" {
			pretty = true
		} else if a == "-s" || a == "--since" {
			since = true
		}
	}

	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		fmt.Fprintln(os.Stderr, "uptime:", err)
		os.Exit(1)
	}
	uptimeSecs := info.Uptime
	bootTime := time.Now().Add(-time.Duration(uptimeSecs) * time.Second)

	if since {
		fmt.Println(bootTime.Format("2006-01-02 15:04:05"))
		return
	}

	hours := uptimeSecs / 3600
	mins := (uptimeSecs % 3600) / 60

	if pretty {
		var parts []string
		days := hours / 24
		hours = hours % 24
		if days > 0 {
			parts = append(parts, fmt.Sprintf("%d day%s", days, pluralS(days)))
		}
		if hours > 0 {
			parts = append(parts, fmt.Sprintf("%d hour%s", hours, pluralS(hours)))
		}
		if mins > 0 {
			parts = append(parts, fmt.Sprintf("%d minute%s", mins, pluralS(mins)))
		}
		fmt.Printf("up %s\n", strings.Join(parts, ", "))
		return
	}

	// Standard output
	now := time.Now().Format("15:04:05")
	loadAvg := float64(info.Loads[0]) / 65536.0
	loadAvg2 := float64(info.Loads[1]) / 65536.0
	loadAvg3 := float64(info.Loads[2]) / 65536.0
	fmt.Printf(" %s up %d:%02d,  %d user%s,  load average: %.2f, %.2f, %.2f\n",
		now, hours, mins, info.Procs, pluralS(int64(info.Procs)), loadAvg, loadAvg2, loadAvg3)
}

