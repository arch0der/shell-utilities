// ss - Socket statistics (Linux /proc/net based)
// Usage: ss [-t] [-u] [-l] [-a] [-n] [-p]
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	tcpFlag  = flag.Bool("t", false, "Show TCP sockets")
	udpFlag  = flag.Bool("u", false, "Show UDP sockets")
	listenF  = flag.Bool("l", false, "Show only listening sockets")
	allFlag  = flag.Bool("a", false, "Show all sockets")
	numericF = flag.Bool("n", false, "Do not resolve service names")
	summaryF = flag.Bool("s", false, "Print summary statistics")
)

type Socket struct {
	proto  string
	local  string
	remote string
	state  string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: ss [-t] [-u] [-l] [-a] [-n] [-s] [filter]")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Default to TCP if nothing specified
	if !*tcpFlag && !*udpFlag {
		*tcpFlag = true
		*udpFlag = true
	}

	if *summaryF {
		printSummary()
		return
	}

	fmt.Printf("%-6s %-12s %-28s %-28s\n", "State", "Recv-Q Send-Q", "Local Address:Port", "Peer Address:Port")

	if *tcpFlag {
		printTCP()
	}
	if *udpFlag {
		printUDP()
	}
}

var tcpStates = map[string]string{
	"01": "ESTABLISHED", "02": "SYN_SENT", "03": "SYN_RECV",
	"04": "FIN_WAIT1", "05": "FIN_WAIT2", "06": "TIME_WAIT",
	"07": "CLOSE", "08": "CLOSE_WAIT", "09": "LAST_ACK",
	"0A": "LISTEN", "0B": "CLOSING",
}

func printTCP() {
	data, err := os.ReadFile("/proc/net/tcp")
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		state := strings.ToUpper(fields[3])
		stateName := tcpStates[state]
		if stateName == "" {
			stateName = state
		}
		if *listenF && stateName != "LISTEN" {
			continue
		}
		if !*allFlag && !*listenF && stateName == "LISTEN" {
			continue
		}
		local := hexToAddr(fields[1])
		remote := hexToAddr(fields[2])
		fmt.Printf("%-6s %-12s %-28s %-28s\n", "tcp", stateName, local, remote)
	}
}

func printUDP() {
	data, err := os.ReadFile("/proc/net/udp")
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		local := hexToAddr(fields[1])
		remote := hexToAddr(fields[2])
		fmt.Printf("%-6s %-12s %-28s %-28s\n", "udp", "UNCONN", local, remote)
	}
}

func printSummary() {
	tcp, _ := os.ReadFile("/proc/net/tcp")
	udp, _ := os.ReadFile("/proc/net/udp")
	tcpLines := len(strings.Split(string(tcp), "\n")) - 2
	udpLines := len(strings.Split(string(udp), "\n")) - 2
	fmt.Println("Socket statistics:")
	fmt.Printf("  TCP: %d\n", tcpLines)
	fmt.Printf("  UDP: %d\n", udpLines)
}

func hexToAddr(h string) string {
	parts := strings.Split(h, ":")
	if len(parts) != 2 {
		return h
	}
	// IP is little-endian hex
	ipHex := parts[0]
	portHex := parts[1]

	var ipBytes [4]byte
	for i := 0; i < 4; i++ {
		var n int
		fmt.Sscanf(ipHex[i*2:i*2+2], "%x", &n)
		ipBytes[3-i] = byte(n)
	}
	ip := net.IP(ipBytes[:])
	var port int
	fmt.Sscanf(portHex, "%x", &port)

	if !*numericF {
		// Could resolve service name here
		_ = port
	}
	return fmt.Sprintf("%s:%d", ip, port)
}
