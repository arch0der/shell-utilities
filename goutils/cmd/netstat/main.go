// netstat - Show open network connections and ports (Linux /proc-based)
// Usage: netstat [-l] [-t] [-u]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	listen = flag.Bool("l", false, "Show only listening sockets")
	tcp    = flag.Bool("t", false, "Show TCP sockets")
	udp    = flag.Bool("u", false, "Show UDP sockets")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: netstat [-l] [-t] [-u]")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Default: show both
	if !*tcp && !*udp {
		*tcp, *udp = true, true
	}

	fmt.Printf("%-6s %-22s %-22s %s\n", "Proto", "Local Address", "Remote Address", "State")

	if *tcp {
		printProc("/proc/net/tcp", "tcp")
	}
	if *udp {
		printProc("/proc/net/udp", "udp")
	}
}

func printProc(path, proto string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "netstat:", err)
		return
	}
	defer f.Close()

	states := map[string]string{
		"01": "ESTABLISHED", "02": "SYN_SENT", "03": "SYN_RECV",
		"04": "FIN_WAIT1", "05": "FIN_WAIT2", "06": "TIME_WAIT",
		"07": "CLOSE", "08": "CLOSE_WAIT", "09": "LAST_ACK",
		"0A": "LISTEN", "0B": "CLOSING",
	}

	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		local := hexToAddr(fields[1])
		remote := hexToAddr(fields[2])
		state := states[strings.ToUpper(fields[3])]

		if *listen && state != "LISTEN" {
			continue
		}

		fmt.Printf("%-6s %-22s %-22s %s\n", proto, local, remote, state)
	}
}

func hexToAddr(h string) string {
	parts := strings.Split(h, ":")
	if len(parts) != 2 {
		return h
	}
	// IP is little-endian hex
	ipHex := parts[0]
	portHex := parts[1]

	b := make([]byte, 4)
	for i := 0; i < 4; i++ {
		v, _ := strconv.ParseUint(ipHex[i*2:i*2+2], 16, 8)
		b[3-i] = byte(v)
	}
	ip := net.IP(b).String()

	port, _ := strconv.ParseInt(portHex, 16, 32)
	return fmt.Sprintf("%s:%d", ip, port)
}
