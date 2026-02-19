// cidr - CIDR network calculator.
//
// Usage:
//
//	cidr NETWORK[/PREFIX]
//	cidr [OPTIONS] NETWORK/PREFIX
//
// Options:
//
//	-l        List all IPs in range (warning: large for /16 and below)
//	-c N      Check if IP N is in the network
//	-s        Split into two subnets
//	-j        JSON output
//
// Examples:
//
//	cidr 192.168.1.0/24
//	cidr 10.0.0.0/8
//	cidr -c 192.168.1.50 192.168.1.0/24    # is IP in range?
//	cidr -l 10.0.0.0/30                     # list all 4 IPs
//	cidr -j 172.16.0.0/12
package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	list    = flag.Bool("l", false, "list all IPs")
	check   = flag.String("c", "", "check if IP is in network")
	split   = flag.Bool("s", false, "split into subnets")
	asJSON  = flag.Bool("j", false, "JSON output")
)

func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

func uint32ToIP(n uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip
}

type Info struct {
	Network   string `json:"network"`
	Mask      string `json:"mask"`
	Broadcast string `json:"broadcast"`
	First     string `json:"first_host"`
	Last      string `json:"last_host"`
	Hosts     uint64 `json:"usable_hosts"`
	Total     uint64 `json:"total_ips"`
	Prefix    int    `json:"prefix_length"`
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: cidr NETWORK/PREFIX")
		os.Exit(1)
	}

	_, network, err := net.ParseCIDR(args[0])
	if err != nil {
		// try adding /32
		_, network, err = net.ParseCIDR(args[0] + "/32")
		if err != nil {
			fmt.Fprintf(os.Stderr, "cidr: invalid network: %v\n", err)
			os.Exit(1)
		}
	}

	prefix, bits := network.Mask.Size()
	if bits != 32 {
		fmt.Fprintln(os.Stderr, "cidr: only IPv4 supported")
		os.Exit(1)
	}

	netIP := ipToUint32(network.IP)
	maskInt := ^uint32(0) << uint(32-prefix)
	broadcast := netIP | ^maskInt
	firstHost := netIP + 1
	lastHost := broadcast - 1
	total := uint64(1) << uint(32-prefix)
	usable := uint64(0)
	if total >= 2 {
		usable = total - 2
	}
	if prefix == 32 {
		firstHost = netIP
		lastHost = netIP
		usable = 1
	}
	if prefix == 31 {
		firstHost = netIP
		lastHost = broadcast
		usable = 2
	}

	maskIP := net.IP(network.Mask).String()

	if *check != "" {
		testIP := net.ParseIP(*check)
		if testIP == nil {
			fmt.Fprintf(os.Stderr, "cidr: invalid IP: %s\n", *check)
			os.Exit(1)
		}
		if network.Contains(testIP) {
			fmt.Printf("%s is IN %s\n", *check, network)
			os.Exit(0)
		}
		fmt.Printf("%s is NOT in %s\n", *check, network)
		os.Exit(1)
	}

	info := Info{
		Network:   network.IP.String(),
		Mask:      maskIP,
		Broadcast: uint32ToIP(broadcast).String(),
		First:     uint32ToIP(firstHost).String(),
		Last:      uint32ToIP(lastHost).String(),
		Hosts:     usable,
		Total:     total,
		Prefix:    prefix,
	}

	if *asJSON {
		b, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(b))
		return
	}

	fmt.Printf("Network:     %s/%d\n", info.Network, prefix)
	fmt.Printf("Mask:        %s\n", info.Mask)
	fmt.Printf("Broadcast:   %s\n", info.Broadcast)
	fmt.Printf("First host:  %s\n", info.First)
	fmt.Printf("Last host:   %s\n", info.Last)
	fmt.Printf("Usable:      %d\n", info.Hosts)
	fmt.Printf("Total IPs:   %d\n", info.Total)

	if *split {
		newPrefix := prefix + 1
		if newPrefix > 32 {
			fmt.Fprintln(os.Stderr, "cidr: cannot split a /32")
			return
		}
		sub1 := fmt.Sprintf("%s/%d", uint32ToIP(netIP), newPrefix)
		mid := netIP + (1 << uint(32-newPrefix))
		sub2 := fmt.Sprintf("%s/%d", uint32ToIP(mid), newPrefix)
		fmt.Printf("\nSubnets:\n  %s\n  %s\n", sub1, sub2)
	}

	if *list {
		if total > 65536 {
			fmt.Fprintf(os.Stderr, "cidr: refusing to list %d IPs (max 65536). Use a larger prefix.\n", total)
			os.Exit(1)
		}
		fmt.Println("\nAll IPs:")
		for i := uint32(0); i < uint32(total); i++ {
			fmt.Println(uint32ToIP(netIP + i))
		}
	}
}
