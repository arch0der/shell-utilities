// xxd - Hex dump of files
// Usage: xxd [-c cols] [-l limit] [-r] [file]
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	cols    = flag.Int("c", 16, "Bytes per line")
	limit   = flag.Int64("l", -1, "Stop after N bytes")
	reverse = flag.Bool("r", false, "Reverse: hex dump back to binary")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: xxd [-c cols] [-l limit] [-r] [file]")
		flag.PrintDefaults()
	}
	flag.Parse()

	var r io.Reader = os.Stdin
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "xxd:", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}

	if *reverse {
		reverseHex(r, os.Stdout)
		return
	}

	buf := make([]byte, *cols)
	offset := int64(0)
	total := int64(0)

	for {
		toRead := int64(*cols)
		if *limit >= 0 && total+toRead > *limit {
			toRead = *limit - total
		}
		if toRead == 0 {
			break
		}
		n, err := r.Read(buf[:toRead])
		if n > 0 {
			// Hex part
			hexStr := fmt.Sprintf("% x", buf[:n])
			// Pad
			padLen := *cols*3 - 1 - len(hexStr)
			if padLen < 0 {
				padLen = 0
			}
			// ASCII part
			ascii := make([]byte, n)
			for i, b := range buf[:n] {
				if b >= 32 && b < 127 {
					ascii[i] = b
				} else {
					ascii[i] = '.'
				}
			}
			fmt.Printf("%08x: %-*s  %s\n", offset, *cols*3-1, hexStr, strings.Repeat(" ", padLen)+string(ascii))
			offset += int64(n)
			total += int64(n)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "xxd:", err)
			os.Exit(1)
		}
	}
}

func reverseHex(r io.Reader, w io.Writer) {
	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xxd:", err)
		os.Exit(1)
	}
	// Parse lines: "offset: hex  ascii"
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Strip offset
		if idx := strings.Index(line, ": "); idx >= 0 {
			line = line[idx+2:]
		}
		// Strip ASCII (after two spaces)
		if idx := strings.Index(line, "  "); idx >= 0 {
			line = line[:idx]
		}
		line = strings.ReplaceAll(line, " ", "")
		b, err := hex.DecodeString(line)
		if err != nil {
			continue
		}
		w.Write(b)
	}
}
