// dig - DNS lookup utility (detailed output)
// Usage: dig [@server] [name] [type]
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: dig [@server] [name] [type]")
		os.Exit(1)
	}

	rtype := "A"
	name := ""

	for _, a := range args {
		if strings.HasPrefix(a, "@") {
			// Custom server (ignored in stdlib, noted only)
			fmt.Fprintf(os.Stderr, "; Using default resolver (custom server %s not supported)\n", a)
		} else if isType(a) {
			rtype = strings.ToUpper(a)
		} else {
			name = a
		}
	}

	if name == "" {
		fmt.Fprintln(os.Stderr, "dig: no name specified")
		os.Exit(1)
	}

	fmt.Printf("; <<>> DiG (Go) <<>> %s %s\n", rtype, name)
	fmt.Printf(";; QUESTION SECTION:\n;%-30s IN %s\n\n", name+".", rtype)
	fmt.Println(";; ANSWER SECTION:")

	start := time.Now()
	switch strings.ToUpper(rtype) {
	case "A", "AAAA":
		addrs, err := net.LookupHost(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dig:", err)
			os.Exit(1)
		}
		for _, a := range addrs {
			fmt.Printf("%-30s 300 IN %-6s %s\n", name+".", rtype, a)
		}
	case "MX":
		records, err := net.LookupMX(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dig:", err)
			os.Exit(1)
		}
		for _, r := range records {
			fmt.Printf("%-30s 300 IN MX    %d %s\n", name+".", r.Pref, r.Host)
		}
	case "NS":
		records, err := net.LookupNS(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dig:", err)
			os.Exit(1)
		}
		for _, r := range records {
			fmt.Printf("%-30s 300 IN NS    %s\n", name+".", r.Host)
		}
	case "TXT":
		records, err := net.LookupTXT(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dig:", err)
			os.Exit(1)
		}
		for _, r := range records {
			fmt.Printf("%-30s 300 IN TXT   %q\n", name+".", r)
		}
	case "CNAME":
		cname, err := net.LookupCNAME(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dig:", err)
			os.Exit(1)
		}
		fmt.Printf("%-30s 300 IN CNAME %s\n", name+".", cname)
	case "PTR":
		names, err := net.LookupAddr(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dig:", err)
			os.Exit(1)
		}
		for _, n := range names {
			fmt.Printf("%-30s 300 IN PTR   %s\n", name+".", n)
		}
	default:
		fmt.Fprintf(os.Stderr, "dig: unsupported type %s\n", rtype)
		os.Exit(1)
	}

	elapsed := time.Since(start)
	fmt.Printf("\n;; Query time: %v\n", elapsed)
	fmt.Printf(";; WHEN: %s\n", time.Now().Format("Mon Jan 02 15:04:05 2006"))
}

func isType(s string) bool {
	types := map[string]bool{"A": true, "AAAA": true, "MX": true, "NS": true,
		"TXT": true, "CNAME": true, "PTR": true, "SOA": true, "SRV": true}
	return types[strings.ToUpper(s)]
}
