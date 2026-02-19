// ifconfig - Show network interface configuration (Linux /proc/net based)
// Usage: ifconfig [interface]
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	filter := ""
	if len(os.Args) > 1 {
		filter = os.Args[1]
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ifconfig:", err)
		os.Exit(1)
	}

	first := true
	for _, iface := range ifaces {
		if filter != "" && iface.Name != filter {
			continue
		}
		if !first {
			fmt.Println()
		}
		first = false

		flags := []string{}
		if iface.Flags&net.FlagUp != 0 {
			flags = append(flags, "UP")
		}
		if iface.Flags&net.FlagBroadcast != 0 {
			flags = append(flags, "BROADCAST")
		}
		if iface.Flags&net.FlagLoopback != 0 {
			flags = append(flags, "LOOPBACK")
		}
		if iface.Flags&net.FlagPointToPoint != 0 {
			flags = append(flags, "POINTTOPOINT")
		}
		if iface.Flags&net.FlagMulticast != 0 {
			flags = append(flags, "MULTICAST")
		}
		if iface.Flags&net.FlagRunning != 0 {
			flags = append(flags, "RUNNING")
		}

		fmt.Printf("%-12s flags=%d<%s>  mtu %d\n",
			iface.Name, iface.Flags, strings.Join(flags, ","), iface.MTU)

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil {
					mask := v.Mask
					fmt.Printf("        inet %s  netmask %d.%d.%d.%d\n",
						v.IP.String(), mask[0], mask[1], mask[2], mask[3])
				} else {
					ones, _ := v.Mask.Size()
					fmt.Printf("        inet6 %s  prefixlen %d\n", v.IP.String(), ones)
				}
			}
		}

		if iface.HardwareAddr != nil {
			fmt.Printf("        ether %s\n", iface.HardwareAddr)
		}

		// Read TX/RX stats from /proc/net/dev
		stats := readIfaceStats(iface.Name)
		if stats != "" {
			fmt.Print(stats)
		}
	}

	if filter != "" && first {
		fmt.Fprintf(os.Stderr, "ifconfig: interface %s not found\n", filter)
		os.Exit(1)
	}
}

func readIfaceStats(name string) string {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, name+":") {
			continue
		}
		parts := strings.Fields(strings.TrimPrefix(line, name+":"))
		if len(parts) < 9 {
			return ""
		}
		return fmt.Sprintf("        RX packets %s  bytes %s\n        TX packets %s  bytes %s\n",
			parts[1], parts[0], parts[9], parts[8])
	}
	return ""
}
