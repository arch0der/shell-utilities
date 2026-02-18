// dns - DNS lookup utility
// Usage: dns [-type A|MX|NS|TXT|CNAME] <host>
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var rtype = flag.String("type", "A", "Record type: A, AAAA, MX, NS, TXT, CNAME")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dns [-type A|AAAA|MX|NS|TXT|CNAME] <host>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	host := flag.Arg(0)

	switch strings.ToUpper(*rtype) {
	case "A", "AAAA":
		addrs, err := net.LookupHost(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dns:", err)
			os.Exit(1)
		}
		for _, a := range addrs {
			fmt.Printf("%s\t%s\t%s\n", host, *rtype, a)
		}
	case "MX":
		records, err := net.LookupMX(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dns:", err)
			os.Exit(1)
		}
		for _, r := range records {
			fmt.Printf("%s\tMX\t%d %s\n", host, r.Pref, r.Host)
		}
	case "NS":
		records, err := net.LookupNS(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dns:", err)
			os.Exit(1)
		}
		for _, r := range records {
			fmt.Printf("%s\tNS\t%s\n", host, r.Host)
		}
	case "TXT":
		records, err := net.LookupTXT(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dns:", err)
			os.Exit(1)
		}
		for _, r := range records {
			fmt.Printf("%s\tTXT\t%q\n", host, r)
		}
	case "CNAME":
		cname, err := net.LookupCNAME(host)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dns:", err)
			os.Exit(1)
		}
		fmt.Printf("%s\tCNAME\t%s\n", host, cname)
	default:
		fmt.Fprintf(os.Stderr, "dns: unsupported type %s\n", *rtype)
		os.Exit(1)
	}
}
