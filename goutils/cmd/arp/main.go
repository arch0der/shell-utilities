// arp - Display and modify the ARP cache (Linux /proc/net/arp)
// Usage: arp [-n] [-a] [-d host] [host]
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	numeric = flag.Bool("n", false, "Show numeric addresses")
	all     = flag.Bool("a", false, "Show all entries")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: arp [-n] [-a] [host]")
		flag.PrintDefaults()
	}
	flag.Parse()

	data, err := os.ReadFile("/proc/net/arp")
	if err != nil {
		fmt.Fprintln(os.Stderr, "arp:", err)
		os.Exit(1)
	}

	filter := ""
	if flag.NArg() > 0 {
		filter = flag.Arg(0)
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return
	}

	// Header
	fmt.Printf("%-18s %-8s %-8s %-20s %s\n", "Address", "HWtype", "HWaddress", "Flags Mask", "Iface")

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		ip := fields[0]
		hwtype := fields[1]
		flags := fields[2]
		mac := fields[3]
		iface := fields[5]

		if filter != "" && ip != filter && iface != filter {
			continue
		}

		hwTypeName := "ether"
		if hwtype == "0x1" {
			hwTypeName = "ether"
		}

		flagName := "C"
		if flags == "0x0" {
			flagName = "I"
		}

		fmt.Printf("%-18s %-8s %-17s %-13s %s\n", ip, hwTypeName, mac, flagName, iface)
	}
}
