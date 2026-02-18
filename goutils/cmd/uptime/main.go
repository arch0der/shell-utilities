// uptime - Show system uptime and load averages
// Usage: uptime
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Read uptime from /proc/uptime
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		fmt.Fprintln(os.Stderr, "uptime:", err)
		os.Exit(1)
	}
	fields := strings.Fields(string(data))
	secs, _ := strconv.ParseFloat(fields[0], 64)
	dur := time.Duration(secs) * time.Second

	days := int(dur.Hours()) / 24
	hours := int(dur.Hours()) % 24
	mins := int(dur.Minutes()) % 60

	// Read load averages from /proc/loadavg
	loadData, _ := os.ReadFile("/proc/loadavg")
	loadFields := strings.Fields(string(loadData))
	load1, load5, load15 := loadFields[0], loadFields[1], loadFields[2]

	now := time.Now().Format("15:04:05")
	uptimeStr := ""
	if days > 0 {
		uptimeStr = fmt.Sprintf("%d day(s), %d:%02d", days, hours, mins)
	} else {
		uptimeStr = fmt.Sprintf("%d:%02d", hours, mins)
	}

	fmt.Printf(" %s  up %s,  load average: %s, %s, %s\n",
		now, uptimeStr, load1, load5, load15)
}
