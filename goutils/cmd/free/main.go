// free - Display memory usage (Linux /proc/meminfo)
// Usage: free [-h] [-m] [-g]
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	human   = flag.Bool("h", false, "Human-readable sizes")
	megs    = flag.Bool("m", false, "Show values in MB")
	gigs    = flag.Bool("g", false, "Show values in GB")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: free [-h] [-m] [-g]")
		flag.PrintDefaults()
	}
	flag.Parse()

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Fprintln(os.Stderr, "free:", err)
		os.Exit(1)
	}

	mem := parseMeminfo(string(data))
	total := mem["MemTotal"]
	free := mem["MemFree"]
	available := mem["MemAvailable"]
	buffers := mem["Buffers"]
	cached := mem["Cached"]
	used := total - free - buffers - cached

	swapTotal := mem["SwapTotal"]
	swapFree := mem["SwapFree"]
	swapUsed := swapTotal - swapFree

	unit := "kB"
	div := int64(1)
	if *megs {
		div = 1024
		unit = "MB"
	} else if *gigs {
		div = 1024 * 1024
		unit = "GB"
	}

	if *human {
		fmt.Printf("%-10s %10s %10s %10s %10s %10s\n",
			"", "total", "used", "free", "buffers", "available")
		fmt.Printf("%-10s %10s %10s %10s %10s %10s\n",
			"Mem:", hfmt(total), hfmt(used), hfmt(free), hfmt(buffers), hfmt(available))
		fmt.Printf("%-10s %10s %10s %10s\n",
			"Swap:", hfmt(swapTotal), hfmt(swapUsed), hfmt(swapFree))
	} else {
		fmt.Printf("%-10s %10s %10s %10s %10s %10s\n",
			"("+unit+")", "total", "used", "free", "buffers", "available")
		fmt.Printf("%-10s %10d %10d %10d %10d %10d\n",
			"Mem:", total/div, used/div, free/div, buffers/div, available/div)
		fmt.Printf("%-10s %10d %10d %10d\n",
			"Swap:", swapTotal/div, swapUsed/div, swapFree/div)
	}
}

func parseMeminfo(s string) map[string]int64 {
	m := map[string]int64{}
	for _, line := range strings.Split(s, "\n") {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSuffix(parts[0], ":")
		val, _ := strconv.ParseInt(parts[1], 10, 64)
		m[key] = val
	}
	return m
}

func hfmt(kb int64) string {
	b := kb * 1024
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1fG", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1fM", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1fK", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%dB", b)
	}
}
