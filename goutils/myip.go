// myip - Print your public IP address(es).
//
// Usage:
//
//	myip [OPTIONS]
//
// Options:
//
//	-4        IPv4 only
//	-6        IPv6 only
//	-l        Also print local/LAN interfaces
//	-j        JSON output
//	-s        Short: just the IP, no labels
//
// Examples:
//
//	myip            # public IPv4
//	myip -l         # all IPs including local
//	myip -j         # JSON output
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	ipv4  = flag.Bool("4", false, "IPv4 only")
	ipv6  = flag.Bool("6", false, "IPv6 only")
	local = flag.Bool("l", false, "include local IPs")
	asJSON = flag.Bool("j", false, "JSON output")
	short = flag.Bool("s", false, "short output")
)

type Result struct {
	Public string            `json:"public,omitempty"`
	Local  map[string]string `json:"local,omitempty"`
}

func fetchPublic(url string) string {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(b))
}

func getPublicIP() string {
	if *ipv6 {
		ip := fetchPublic("https://api6.ipify.org")
		if ip != "" {
			return ip
		}
		return fetchPublic("https://ipv6.icanhazip.com")
	}
	ip := fetchPublic("https://api4.ipify.org")
	if ip != "" {
		return ip
	}
	return fetchPublic("https://ipv4.icanhazip.com")
}

func getLocalIPs() map[string]string {
	result := make(map[string]string)
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			if *ipv4 && ip.To4() == nil {
				continue
			}
			if *ipv6 && ip.To4() != nil {
				continue
			}
			result[iface.Name] = ip.String()
		}
	}
	return result
}

func main() {
	flag.Parse()

	pub := getPublicIP()
	locals := map[string]string{}
	if *local {
		locals = getLocalIPs()
	}

	if *asJSON {
		r := Result{Public: pub, Local: locals}
		b, _ := json.MarshalIndent(r, "", "  ")
		fmt.Println(string(b))
		return
	}

	if *short {
		if pub != "" {
			fmt.Println(pub)
		}
		for _, ip := range locals {
			fmt.Println(ip)
		}
		return
	}

	if pub != "" {
		fmt.Printf("public: %s\n", pub)
	} else {
		fmt.Fprintln(os.Stderr, "myip: could not determine public IP")
	}
	for iface, ip := range locals {
		fmt.Printf("%-12s %s\n", iface+":", ip)
	}
}
