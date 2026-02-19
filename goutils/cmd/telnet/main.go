// telnet - Simple TCP connection tool (telnet protocol subset)
// Usage: telnet [host [port]]
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		// Interactive mode: prompt for host/port
		fmt.Print("telnet> ")
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if strings.HasPrefix(line, "open ") {
				parts := strings.Fields(line[5:])
				if len(parts) >= 2 {
					connect(parts[0], parts[1])
				} else if len(parts) == 1 {
					connect(parts[0], "23")
				}
			} else if line == "quit" || line == "exit" {
				break
			} else if line != "" {
				fmt.Println("Commands: open <host> [port], quit")
			}
			fmt.Print("telnet> ")
		}
		return
	}

	host := os.Args[1]
	port := "23"
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	connect(host, port)
}

func connect(host, port string) {
	fmt.Printf("Trying %s...\n", host)
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "telnet:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to %s.\nEscape character is '^]'.\n", host)

	// Negotiate telnet options: respond to IAC commands
	done := make(chan struct{})

	// Read from server -> stdout
	go func() {
		defer close(done)
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				// Strip telnet IAC sequences
				out := stripTelnet(buf[:n])
				os.Stdout.Write(out)
			}
			if err != nil {
				if err != io.EOF {
					fmt.Fprintln(os.Stderr)
				}
				break
			}
		}
		fmt.Println("\nConnection closed by foreign host.")
	}()

	// Read from stdin -> server
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			line := sc.Text() + "\r\n"
			conn.Write([]byte(line))
		}
		conn.Close()
	}()

	<-done
}

// Strip telnet IAC (Interpret As Command) sequences
func stripTelnet(data []byte) []byte {
	const IAC = 0xFF
	out := make([]byte, 0, len(data))
	for i := 0; i < len(data); i++ {
		if data[i] == IAC && i+1 < len(data) {
			cmd := data[i+1]
			if cmd == IAC {
				out = append(out, IAC)
				i++
			} else if cmd >= 0xFB && cmd <= 0xFE && i+2 < len(data) {
				// DO/DONT/WILL/WONT + option: skip 3 bytes
				i += 2
			} else {
				i++
			}
		} else {
			out = append(out, data[i])
		}
	}
	return out
}
