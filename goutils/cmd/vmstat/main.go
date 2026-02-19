// vmstat - Virtual memory statistics (Linux /proc based)
// Usage: vmstat [-a] [-s] [delay [count]]
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
	all     = flag.Bool("a", false, "Show active/inactive memory")
	summary = flag.Bool("s", false, "Show event counters and memory stats")
)

type MemInfo map[string]int64

func readMemInfo() MemInfo {
	data, _ := os.ReadFile("/proc/meminfo")
	m := MemInfo{}
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			key := strings.TrimSuffix(parts[0], ":")
			val, _ := strconv.ParseInt(parts[1], 10, 64)
			m[key] = val
		}
	}
	return m
}

type StatInfo map[string]int64

func readStat() StatInfo {
	data, _ := os.ReadFile("/proc/stat")
	s := StatInfo{}
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		val, err := strconv.ParseInt(parts[1], 10, 64)
		if err == nil {
			s[parts[0]] = val
		}
	}
	return s
}

func printHeader(active bool) {
	if active {
		fmt.Printf("procs --------memory (kB)-------- ---swap-- -----io---- -system-- ------cpu-----\n")
		fmt.Printf(" r  b   swpd   free   inact  active   si   so    bi    bo   in   cs us sy id wa st\n")
	} else {
		fmt.Printf("procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----\n")
		fmt.Printf(" r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st\n")
	}
}

func printStats(mem MemInfo, stat StatInfo, active bool) {
	free := mem["MemFree"]
	buff := mem["Buffers"]
	cache := mem["Cached"]
	swpd := mem["SwapTotal"] - mem["SwapFree"]
	inact := mem["Inactive"]
	act := mem["Active"]

	// CPU stats from /proc/stat: cpu user nice system idle iowait irq softirq
	cpuFields := []int64{}
	data, _ := os.ReadFile("/proc/stat")
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "cpu ") {
			for _, f := range strings.Fields(line)[1:] {
				n, _ := strconv.ParseInt(f, 10, 64)
				cpuFields = append(cpuFields, n)
			}
			break
		}
	}
	total := int64(0)
	for _, v := range cpuFields {
		total += v
	}
	if total == 0 {
		total = 1
	}

	us, sy, id, wa := int64(0), int64(0), int64(0), int64(0)
	if len(cpuFields) >= 4 {
		us = cpuFields[0] * 100 / total
		sy = cpuFields[2] * 100 / total
		id = cpuFields[3] * 100 / total
	}
	if len(cpuFields) >= 5 {
		wa = cpuFields[4] * 100 / total
	}

	if active {
		fmt.Printf(" 0  0 %6d %6d %6d %6d    0    0     0     0    0    0 %2d %2d %2d %2d  0\n",
			swpd, free, inact, act, us, sy, id, wa)
	} else {
		fmt.Printf(" 0  0 %6d %6d %6d %6d    0    0     0     0    0    0 %2d %2d %2d %2d  0\n",
			swpd, free, buff, cache, us, sy, id, wa)
	}
}

func printSummary(mem MemInfo) {
	fmt.Printf("%12d K total memory\n", mem["MemTotal"])
	fmt.Printf("%12d K used memory\n", mem["MemTotal"]-mem["MemFree"]-mem["Buffers"]-mem["Cached"])
	fmt.Printf("%12d K active memory\n", mem["Active"])
	fmt.Printf("%12d K inactive memory\n", mem["Inactive"])
	fmt.Printf("%12d K free memory\n", mem["MemFree"])
	fmt.Printf("%12d K buffer memory\n", mem["Buffers"])
	fmt.Printf("%12d K swap cache\n", mem["Cached"])
	fmt.Printf("%12d K total swap\n", mem["SwapTotal"])
	fmt.Printf("%12d K used swap\n", mem["SwapTotal"]-mem["SwapFree"])
	fmt.Printf("%12d K free swap\n", mem["SwapFree"])
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: vmstat [-a] [-s] [delay [count]]")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	delay := 0
	count := 1
	if len(args) > 0 {
		delay, _ = strconv.Atoi(args[0])
	}
	if len(args) > 1 {
		count, _ = strconv.Atoi(args[1])
		if delay == 0 {
			delay = 1
		}
	}

	mem := readMemInfo()
	stat := readStat()

	if *summary {
		printSummary(mem)
		return
	}

	printHeader(*all)
	for i := 0; i < count; i++ {
		if i > 0 {
			time.Sleep(time.Duration(delay) * time.Second)
			mem = readMemInfo()
			stat = readStat()
		}
		printStats(mem, stat, *all)
	}
}
