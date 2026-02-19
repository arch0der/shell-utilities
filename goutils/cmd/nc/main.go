// nc - Netcat: read/write TCP or UDP connections
// Usage: nc [-l] [-u] [-p port] [host] [port]
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	listen  = flag.Bool("l", false, "Listen mode")
	udp     = flag.Bool("u", false, "Use UDP instead of TCP")
	port    = flag.String("p", "", "Local port (listen mode)")
	verbose = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: nc [-l] [-u] [-v] [-p port] [host] [port]")
		flag.PrintDefaults()
	}
	flag.Parse()

	proto := "tcp"
	if *udp {
		proto = "udp"
	}

	if *listen {
		addr := ":" + *port
		if flag.NArg() > 0 {
			addr = flag.Arg(0)
			if flag.NArg() > 1 {
				addr = flag.Arg(0) + ":" + flag.Arg(1)
			}
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "Listening on %s (%s)\n", addr, proto)
		}
		ln, err := net.Listen(proto, addr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "nc:", err)
			os.Exit(1)
		}
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "nc:", err)
			os.Exit(1)
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "Connection from %s\n", conn.RemoteAddr())
		}
		go io.Copy(conn, os.Stdin)
		io.Copy(os.Stdout, conn)
		conn.Close()
	} else {
		if flag.NArg() < 2 {
			flag.Usage()
			os.Exit(1)
		}
		host, p := flag.Arg(0), flag.Arg(1)
		addr := host + ":" + p
		if *verbose {
			fmt.Fprintf(os.Stderr, "Connecting to %s (%s)\n", addr, proto)
		}
		conn, err := net.Dial(proto, addr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "nc:", err)
			os.Exit(1)
		}
		defer conn.Close()
		go io.Copy(conn, os.Stdin)
		io.Copy(os.Stdout, conn)
	}
}
