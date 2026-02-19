// iprange - expand CIDR or IP ranges, check membership, list IPs
package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

func ip2int(ip net.IP) uint32 { return binary.BigEndian.Uint32(ip.To4()) }
func int2ip(n uint32) net.IP  { ip := make(net.IP, 4); binary.BigEndian.PutUint32(ip, n); return ip }

func expandCIDR(cidr string) ([]string, error) {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil { return nil, err }
	start := ip2int(network.IP)
	mask := binary.BigEndian.Uint32(network.Mask)
	end := start | ^mask
	var ips []string
	for i := start; i <= end; i++ { ips = append(ips, int2ip(i).String()) }
	return ips, nil
}

func expandRange(from, to string) ([]string, error) {
	a := net.ParseIP(from).To4()
	b := net.ParseIP(to).To4()
	if a == nil || b == nil { return nil, fmt.Errorf("invalid IP") }
	s, e := ip2int(a), ip2int(b)
	if s > e { s, e = e, s }
	if e-s > 65536 { return nil, fmt.Errorf("range too large (>65536)") }
	var ips []string
	for i := s; i <= e; i++ { ips = append(ips, int2ip(i).String()) }
	return ips, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  iprange expand <CIDR>                  list all IPs in CIDR
  iprange expand <start_ip> <end_ip>    list IPs in range
  iprange count  <CIDR>                  print count only
  iprange check  <ip> <CIDR>             check if IP is in CIDR`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 { usage() }
	cmd := os.Args[1]
	switch cmd {
	case "expand":
		if strings.Contains(os.Args[2], "/") {
			ips, err := expandCIDR(os.Args[2])
			if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
			for _, ip := range ips { fmt.Println(ip) }
		} else {
			if len(os.Args) < 4 { usage() }
			ips, err := expandRange(os.Args[2], os.Args[3])
			if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
			for _, ip := range ips { fmt.Println(ip) }
		}
	case "count":
		ips, err := expandCIDR(os.Args[2])
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		fmt.Println(len(ips))
	case "check":
		if len(os.Args) < 4 { usage() }
		ip := net.ParseIP(os.Args[2])
		_, network, err := net.ParseCIDR(os.Args[3])
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		if network.Contains(ip) {
			fmt.Printf("%s IS in %s\n", ip, os.Args[3])
		} else {
			fmt.Printf("%s is NOT in %s\n", ip, os.Args[3]); os.Exit(1)
		}
	default:
		usage()
	}
}
