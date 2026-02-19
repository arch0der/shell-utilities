// host - DNS lookup utility
// Usage: host [-t type] [-a] <name> [server]
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	rtype = flag.String("t", "", "Query type: A, AAAA, MX, NS, TXT, CNAME, PTR")
	all   = flag.Bool("a", false, "Query all record types")
	short = flag.Bool("s", false, "Short output (just the answer)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: host [-t type] [-a] [-s] <name> [server]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	name := flag.Arg(0)
	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "; Note: custom DNS server %s not supported, using system resolver\n", flag.Arg(1))
	}

	types := []string{*rtype}
	if *all || *rtype == "" {
		// Default: try A, AAAA, MX
		if *all {
			types = []string{"A", "AAAA", "MX", "NS", "TXT", "CNAME"}
		} else {
			types = []string{"A", "AAAA", "MX"}
		}
	}

	found := false
	for _, t := range types {
		switch strings.ToUpper(t) {
		case "A", "AAAA", "":
			addrs, err := net.LookupHost(name)
			if err != nil {
				continue
			}
			for _, a := range addrs {
				found = true
				if *short {
					fmt.Println(a)
				} else {
					fmt.Printf("%s has address %s\n", name, a)
				}
			}
		case "MX":
			records, err := net.LookupMX(name)
			if err != nil {
				continue
			}
			for _, r := range records {
				found = true
				if *short {
					fmt.Printf("%d %s\n", r.Pref, r.Host)
				} else {
					fmt.Printf("%s mail is handled by %d %s\n", name, r.Pref, r.Host)
				}
			}
		case "NS":
			records, err := net.LookupNS(name)
			if err != nil {
				continue
			}
			for _, r := range records {
				found = true
				if *short {
					fmt.Println(r.Host)
				} else {
					fmt.Printf("%s name server %s\n", name, r.Host)
				}
			}
		case "TXT":
			records, err := net.LookupTXT(name)
			if err != nil {
				continue
			}
			for _, r := range records {
				found = true
				if *short {
					fmt.Println(r)
				} else {
					fmt.Printf("%s descriptive text \"%s\"\n", name, r)
				}
			}
		case "CNAME":
			cname, err := net.LookupCNAME(name)
			if err != nil {
				continue
			}
			found = true
			if *short {
				fmt.Println(cname)
			} else {
				fmt.Printf("%s is an alias for %s\n", name, cname)
			}
		case "PTR":
			names, err := net.LookupAddr(name)
			if err != nil {
				continue
			}
			for _, n := range names {
				found = true
				if *short {
					fmt.Println(n)
				} else {
					fmt.Printf("%s domain name pointer %s\n", name, n)
				}
			}
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Host %s not found\n", name)
		os.Exit(1)
	}
}
