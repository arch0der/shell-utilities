// ip - Show/manipulate routing, devices, and tunnels (Linux, subset)
// Usage: ip [addr|link|route|neigh] [show|add|del] [args]
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcmd := strings.ToLower(os.Args[1])
	action := "show"
	if len(os.Args) > 2 {
		action = strings.ToLower(os.Args[2])
	}

	switch subcmd {
	case "addr", "address", "a":
		ipAddr(action, os.Args[3:])
	case "link", "l":
		ipLink(action, os.Args[3:])
	case "route", "r":
		ipRoute(action, os.Args[3:])
	case "neigh", "n":
		ipNeigh(action, os.Args[3:])
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "ip: unknown command %q\n", subcmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `Usage: ip OBJECT COMMAND [args]
Objects: addr (a), link (l), route (r), neigh (n)
Commands: show, add, del, set`)
}

func ipAddr(action string, args []string) {
	switch action {
	case "show", "list", "":
		ifaces, err := net.Interfaces()
		if err != nil {
			fmt.Fprintln(os.Stderr, "ip addr:", err)
			return
		}
		for _, iface := range ifaces {
			fmt.Printf("%d: %s: <", iface.Index, iface.Name)
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
			if iface.Flags&net.FlagMulticast != 0 {
				flags = append(flags, "MULTICAST")
			}
			fmt.Printf("%s> mtu %d\n", strings.Join(flags, ","), iface.MTU)
			if iface.HardwareAddr != nil {
				fmt.Printf("    link/ether %s\n", iface.HardwareAddr)
			}
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					if v.IP.To4() != nil {
						ones, _ := v.Mask.Size()
						fmt.Printf("    inet %s/%d\n", v.IP, ones)
					} else {
						ones, _ := v.Mask.Size()
						fmt.Printf("    inet6 %s/%d\n", v.IP, ones)
					}
				}
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "ip addr: unsupported action %q (show only)\n", action)
	}
}

func ipLink(action string, args []string) {
	switch action {
	case "show", "list", "":
		ifaces, err := net.Interfaces()
		if err != nil {
			fmt.Fprintln(os.Stderr, "ip link:", err)
			return
		}
		for _, iface := range ifaces {
			state := "DOWN"
			if iface.Flags&net.FlagUp != 0 {
				state = "UP"
			}
			fmt.Printf("%d: %s: state %s mtu %d\n", iface.Index, iface.Name, state, iface.MTU)
			mac := "00:00:00:00:00:00"
			if iface.HardwareAddr != nil {
				mac = iface.HardwareAddr.String()
			}
			fmt.Printf("    link/ether %s\n", mac)
		}
	default:
		fmt.Fprintf(os.Stderr, "ip link: unsupported action %q (show only)\n", action)
	}
}

func ipRoute(action string, args []string) {
	switch action {
	case "show", "list", "":
		// Read from /proc/net/route
		data, err := os.ReadFile("/proc/net/route")
		if err != nil {
			fmt.Fprintln(os.Stderr, "ip route:", err)
			return
		}
		lines := strings.Split(string(data), "\n")
		fmt.Printf("%-20s %-16s %-16s %-6s %s\n", "Destination", "Gateway", "Genmask", "Flags", "Iface")
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			if len(fields) < 8 {
				continue
			}
			iface := fields[0]
			dest := hexToIP(fields[1])
			gw := hexToIP(fields[2])
			mask := hexToIP(fields[7])
			flags := parseRouteFlags(fields[3])
			fmt.Printf("%-20s %-16s %-16s %-6s %s\n", dest, gw, mask, flags, iface)
		}
	default:
		fmt.Fprintf(os.Stderr, "ip route: unsupported action %q\n", action)
	}
}

func ipNeigh(action string, args []string) {
	switch action {
	case "show", "list", "":
		// Read ARP table from /proc/net/arp
		data, err := os.ReadFile("/proc/net/arp")
		if err != nil {
			fmt.Fprintln(os.Stderr, "ip neigh:", err)
			return
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			if len(fields) < 6 {
				continue
			}
			ip := fields[0]
			mac := fields[3]
			iface := fields[5]
			fmt.Printf("%s dev %s lladdr %s REACHABLE\n", ip, iface, mac)
		}
	default:
		fmt.Fprintf(os.Stderr, "ip neigh: unsupported action %q\n", action)
	}
}

func hexToIP(h string) string {
	// Little-endian hex IP
	if len(h) != 8 {
		return h
	}
	parts := make([]byte, 4)
	for i := 0; i < 4; i++ {
		b, _ := fmt.Sscanf(h[i*2:i*2+2], "%x", new(int))
		_ = b
		var n int
		fmt.Sscanf(h[i*2:i*2+2], "%x", &n)
		parts[3-i] = byte(n)
	}
	return net.IP(parts).String()
}

func parseRouteFlags(hexFlags string) string {
	var n int
	fmt.Sscanf(hexFlags, "%x", &n)
	flags := ""
	if n&0x1 != 0 {
		flags += "U"
	}
	if n&0x2 != 0 {
		flags += "G"
	}
	if n&0x4 != 0 {
		flags += "H"
	}
	return flags
}
