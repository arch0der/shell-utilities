// nslookup - Query DNS servers
// Usage: nslookup [-type=TYPE] <name> [server]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var rtype = flag.String("type", "A", "Query type: A, AAAA, MX, NS, TXT, CNAME, PTR")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: nslookup [-type=TYPE] <name> [server]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		// Interactive mode
		runInteractive()
		return
	}

	name := flag.Arg(0)
	server := ""
	if flag.NArg() > 1 {
		server = flag.Arg(1)
		fmt.Fprintf(os.Stderr, "Server:\t\t%s\n", server)
		fmt.Fprintf(os.Stderr, "(Note: custom server not supported, using system resolver)\n\n")
	}

	doLookup(name, strings.ToUpper(*rtype))
}

func doLookup(name, qtype string) {
	fmt.Printf("Server:\t\t(system resolver)\n")
	fmt.Printf("Address:\tsystem\n\n")
	fmt.Printf("Name:\t%s\n", name)

	switch qtype {
	case "A", "AAAA", "":
		addrs, err := net.LookupHost(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "nslookup:", err)
			return
		}
		for _, a := range addrs {
			fmt.Printf("Address:\t%s\n", a)
		}
	case "MX":
		records, err := net.LookupMX(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "nslookup:", err)
			return
		}
		fmt.Println()
		for _, r := range records {
			fmt.Printf("%s\tmail exchanger = %d %s\n", name, r.Pref, r.Host)
		}
	case "NS":
		records, err := net.LookupNS(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "nslookup:", err)
			return
		}
		fmt.Println()
		for _, r := range records {
			fmt.Printf("%s\tnameserver = %s\n", name, r.Host)
		}
	case "TXT":
		records, err := net.LookupTXT(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "nslookup:", err)
			return
		}
		fmt.Println()
		for _, r := range records {
			fmt.Printf("%s\ttext = \"%s\"\n", name, r)
		}
	case "CNAME":
		cname, err := net.LookupCNAME(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "nslookup:", err)
			return
		}
		fmt.Printf("%s\tcanonical name = %s\n", name, cname)
	case "PTR":
		names, err := net.LookupAddr(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "nslookup:", err)
			return
		}
		for _, n := range names {
			fmt.Printf("%s\tpointer = %s\n", name, n)
		}
	default:
		fmt.Fprintf(os.Stderr, "nslookup: unsupported type %s\n", qtype)
	}
}

func runInteractive() {
	fmt.Println("nslookup (type 'exit' or 'quit' to quit, 'set type=X' to change query type)")
	sc := bufio.NewScanner(os.Stdin)
	qtype := "A"
	for {
		fmt.Print("> ")
		if !sc.Scan() {
			break
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "set type=") {
			qtype = strings.ToUpper(strings.TrimPrefix(strings.ToLower(line), "set type="))
			fmt.Printf("Query type set to %s\n", qtype)
			continue
		}
		doLookup(line, qtype)
		fmt.Println()
	}
}
