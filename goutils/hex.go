// hex - Hex encode or decode data.
//
// Usage:
//
//	hex [OPTIONS] [FILE...]
//	echo "hello" | hex
//	echo "68656c6c6f" | hex -d
//
// Options:
//
//	-d        Decode hex to binary
//	-u        Uppercase hex output
//	-g N      Group bytes with space every N bytes (default: 0)
//	-w N      Wrap output at N bytes per line (default: 0)
//	-x        Canonical xxd-style dump (offset + hex + ascii)
//
// Examples:
//
//	echo "hello" | hex              # 68656c6c6f0a
//	echo "68656c6c6f" | hex -d      # hello
//	hex -x file.bin                 # xxd-style dump
//	hex -g 4 -w 16 file.bin         # grouped, 16 bytes/line
package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	decode  = flag.Bool("d", false, "decode")
	upper   = flag.Bool("u", false, "uppercase")
	group   = flag.Int("g", 0, "group every N bytes")
	width   = flag.Int("w", 0, "bytes per line")
	xxd     = flag.Bool("x", false, "xxd-style dump")
)

func encodeHex(data []byte) string {
	s := hex.EncodeToString(data)
	if *upper {
		s = strings.ToUpper(s)
	}
	return s
}

func xxdDump(r io.Reader) {
	buf := make([]byte, 16)
	offset := 0
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for {
		n, err := r.Read(buf)
		if n > 0 {
			// hex part
			hexPart := ""
			for i := 0; i < 16; i++ {
				if i == 8 {
					hexPart += " "
				}
				if i < n {
					hexPart += fmt.Sprintf("%02x ", buf[i])
				} else {
					hexPart += "   "
				}
			}
			// ascii part
			ascii := ""
			for i := 0; i < n; i++ {
				c := buf[i]
				if c >= 32 && c < 127 {
					ascii += string(c)
				} else {
					ascii += "."
				}
			}
			fmt.Fprintf(w, "%08x: %s |%s|\n", offset, hexPart, ascii)
			offset += n
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "hex: %v\n", err)
			break
		}
	}
}

func processEncode(r io.Reader) {
	if *xxd {
		xxdDump(r)
		return
	}

	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hex: %v\n", err)
		os.Exit(1)
	}

	encoded := encodeHex(data)

	if *group == 0 && *width == 0 {
		fmt.Println(encoded)
		return
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	// encoded is pairs of hex chars; 1 byte = 2 chars
	col := 0
	for i := 0; i < len(encoded); i += 2 {
		if col > 0 {
			if *group > 0 && col%*group == 0 {
				w.WriteString(" ")
			}
			if *width > 0 && col%*width == 0 {
				w.WriteString("\n")
				col = 0
			}
		}
		w.WriteString(encoded[i : i+2])
		col++
	}
	w.WriteString("\n")
}

func processDecode(r io.Reader) {
	sc := bufio.NewScanner(r)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for sc.Scan() {
		line := strings.ReplaceAll(sc.Text(), " ", "")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		b, err := hex.DecodeString(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hex: %v\n", err)
			os.Exit(1)
		}
		w.Write(b)
	}
}

func main() {
	flag.Parse()
	files := flag.Args()

	process := processEncode
	if *decode {
		process = processDecode
	}

	if len(files) == 0 {
		process(os.Stdin)
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hex: %v\n", err)
			os.Exit(1)
		}
		process(fh)
		fh.Close()
	}
}
