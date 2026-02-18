// ping - ICMP ping a host (requires root or CAP_NET_RAW)
// Usage: ping [-c count] [-i interval] <host>
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

var (
	count    = flag.Int("c", 4, "Number of pings (0 = infinite)")
	interval = flag.Float64("i", 1.0, "Interval between pings (seconds)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: ping [-c count] [-i interval] <host>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	host := flag.Arg(0)

	addrs, err := net.LookupHost(host)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ping:", err)
		os.Exit(1)
	}
	ip := addrs[0]
	fmt.Printf("PING %s (%s)\n", host, ip)

	conn, err := net.Dial("ip4:icmp", ip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ping: %v (try running with sudo)\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	id := uint16(rand.Intn(0xffff))
	sent, recv := 0, 0

	for i := 0; *count == 0 || i < *count; i++ {
		seq := uint16(i)
		pkt := makeICMP(id, seq)
		start := time.Now()
		conn.SetDeadline(time.Now().Add(3 * time.Second))
		_, err := conn.Write(pkt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ping: write:", err)
			break
		}
		sent++
		buf := make([]byte, 1500)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Request timeout")
		} else {
			rtt := time.Since(start)
			ttl := buf[8] // TTL is at offset 8 in IP header
			_ = n
			fmt.Printf("64 bytes from %s: icmp_seq=%d ttl=%d time=%v\n",
				ip, seq, ttl, rtt.Round(time.Microsecond))
			recv++
		}
		if *count == 0 || i < *count-1 {
			time.Sleep(time.Duration(*interval * float64(time.Second)))
		}
	}

	loss := 100.0
	if sent > 0 {
		loss = float64(sent-recv) / float64(sent) * 100
	}
	fmt.Printf("\n--- %s ping statistics ---\n", host)
	fmt.Printf("%d packets transmitted, %d received, %.0f%% packet loss\n", sent, recv, loss)
}

func makeICMP(id, seq uint16) []byte {
	pkt := make([]byte, 8)
	pkt[0] = 8 // ICMP Echo Request
	pkt[1] = 0
	binary.BigEndian.PutUint16(pkt[4:], id)
	binary.BigEndian.PutUint16(pkt[6:], seq)
	// Checksum
	cs := checksum(pkt)
	binary.BigEndian.PutUint16(pkt[2:], cs)
	return pkt
}

func checksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for sum>>16 != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return ^uint16(sum)
}
