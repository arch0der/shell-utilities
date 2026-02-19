// port - Check if a TCP/UDP port is open on a host.
//
// Usage:
//
//	port [OPTIONS] HOST PORT
//	port [OPTIONS] HOST:PORT
//
// Options:
//
//	-u          Use UDP instead of TCP
//	-t TIMEOUT  Connection timeout (default: 3s)
//	-q          Quiet: exit code only (0=open, 1=closed)
//	-w          Wait mode: retry until open (use with -t for total timeout)
//	-i INTERVAL Retry interval in wait mode (default: 1s)
//
// Examples:
//
//	port google.com 443              # check HTTPS
//	port -q localhost 5432           # quiet, exit code only
//	port -w -t 30s localhost 8080    # wait up to 30s for port to open
//	port -u 8.8.8.8 53               # UDP DNS check
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	udp      = flag.Bool("u", false, "UDP")
	timeout  = flag.Duration("t", 3*time.Second, "timeout")
	quiet    = flag.Bool("q", false, "quiet")
	wait     = flag.Bool("w", false, "wait until open")
	interval = flag.Duration("i", time.Second, "retry interval")
)

func check(network, addr string, to time.Duration) bool {
	if network == "udp" {
		// UDP: just resolve and attempt; no real handshake
		conn, err := net.DialTimeout(network, addr, to)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}
	conn, err := net.DialTimeout(network, addr, to)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func main() {
	flag.Parse()
	args := flag.Args()

	var host, port string
	switch len(args) {
	case 1:
		h, p, err := net.SplitHostPort(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "port: usage: port HOST PORT or HOST:PORT")
			os.Exit(2)
		}
		host, port = h, p
	case 2:
		host, port = args[0], args[1]
	default:
		fmt.Fprintln(os.Stderr, "port: usage: port HOST PORT")
		os.Exit(2)
	}

	network := "tcp"
	if *udp {
		network = "udp"
	}
	addr := net.JoinHostPort(host, port)

	if *wait {
		deadline := time.Now().Add(*timeout)
		for {
			if check(network, addr, 2*time.Second) {
				if !*quiet {
					fmt.Printf("%s is open\n", addr)
				}
				os.Exit(0)
			}
			if time.Now().After(deadline) {
				break
			}
			if !*quiet {
				fmt.Printf("waiting for %s...\n", addr)
			}
			time.Sleep(*interval)
		}
		if !*quiet {
			fmt.Printf("%s did not open within %v\n", addr, *timeout)
		}
		os.Exit(1)
	}

	if check(network, addr, *timeout) {
		if !*quiet {
			fmt.Printf("%s is open\n", addr)
		}
		os.Exit(0)
	}
	if !*quiet {
		fmt.Printf("%s is closed\n", addr)
	}
	os.Exit(1)
}
