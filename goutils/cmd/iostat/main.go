// iostat - Report CPU and disk I/O statistics (Linux /proc based)
// Usage: iostat [-x] [-d] [-c] [interval [count]]
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	extended = flag.Bool("x", false, "Extended statistics")
	diskOnly = flag.Bool("d", false, "Disk statistics only")
	cpuOnly  = flag.Bool("c", false, "CPU statistics only")
)

type DiskStat struct {
	name                        string
	readsCompleted, writesCompleted uint64
	readBytes, writeBytes       uint64
	readTime, writeTime         uint64
}

type CPUStat struct {
	user, nice, system, idle, iowait, irq, softirq uint64
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: iostat [-x] [-d] [-c] [interval [count]]")
		flag.PrintDefaults()
	}
	flag.Parse()

	interval := 1
	count := 1
	args := flag.Args()
	if len(args) > 0 {
		interval, _ = strconv.Atoi(args[0])
	}
	if len(args) > 1 {
		count, _ = strconv.Atoi(args[1])
	}

	fmt.Printf("Linux iostat (Go)\n\n")

	for i := 0; i < count; i++ {
		if i > 0 {
			time.Sleep(time.Duration(interval) * time.Second)
		}
		fmt.Printf("%s\n", time.Now().Format("01/02/2006 03:04:05 PM"))

		if !*diskOnly {
			printCPU()
		}
		if !*cpuOnly {
			printDisk(*extended)
		}
		fmt.Println()
	}
}

func printCPU() {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "cpu ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}
		vals := make([]float64, len(fields)-1)
		total := 0.0
		for i, f := range fields[1:] {
			n, _ := strconv.ParseFloat(f, 64)
			vals[i] = n
			total += n
		}
		if total == 0 {
			total = 1
		}
		fmt.Printf("avg-cpu:  %%user   %%nice %%system %%iowait  %%steal   %%idle\n")
		fmt.Printf("         %6.2f  %6.2f  %6.2f  %6.2f  %6.2f  %6.2f\n\n",
			vals[0]/total*100, vals[1]/total*100, vals[2]/total*100,
			vals[3]/total*100, 0.0, vals[4]/total*100)
		break
	}
}

func printDisk(extended bool) {
	data, err := os.ReadFile("/proc/diskstats")
	if err != nil {
		return
	}

	if extended {
		fmt.Printf("%-10s %8s %8s %12s %12s %8s %8s\n",
			"Device", "r/s", "w/s", "rkB/s", "wkB/s", "r_await", "w_await")
	} else {
		fmt.Printf("%-10s %8s %8s %12s %12s\n",
			"Device", "tps", "kB_read/s", "kB_wrtn/s", "kB_dscd/s")
	}

	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 14 {
			continue
		}
		name := fields[2]
		// Skip loop and ram devices
		if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") {
			continue
		}
		reads, _ := strconv.ParseFloat(fields[3], 64)
		writes, _ := strconv.ParseFloat(fields[7], 64)
		readSectors, _ := strconv.ParseFloat(fields[5], 64)
		writeSectors, _ := strconv.ParseFloat(fields[9], 64)
		readTime, _ := strconv.ParseFloat(fields[6], 64)
		writeTime, _ := strconv.ParseFloat(fields[10], 64)

		readKB := readSectors / 2  // 512-byte sectors to KB
		writeKB := writeSectors / 2

		if extended {
			rAwait := 0.0
			if reads > 0 {
				rAwait = readTime / reads
			}
			wAwait := 0.0
			if writes > 0 {
				wAwait = writeTime / writes
			}
			fmt.Printf("%-10s %8.2f %8.2f %12.2f %12.2f %8.2f %8.2f\n",
				name, reads, writes, readKB, writeKB, rAwait, wAwait)
		} else {
			tps := reads + writes
			fmt.Printf("%-10s %8.2f %12.2f %12.2f %12.2f\n",
				name, tps, readKB, writeKB, 0.0)
		}
	}
}
