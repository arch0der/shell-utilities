// genhex - hex dump, encode/decode, and random hex generation
package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func hexDump(data []byte) string {
	var sb strings.Builder
	for i := 0; i < len(data); i += 16 {
		end := i + 16; if end > len(data) { end = len(data) }
		chunk := data[i:end]
		sb.WriteString(fmt.Sprintf("%08x  ", i))
		for j, b := range chunk {
			sb.WriteString(fmt.Sprintf("%02x ", b))
			if j == 7 { sb.WriteString(" ") }
		}
		padding := 16 - len(chunk)
		sb.WriteString(strings.Repeat("   ", padding))
		if padding > 8 { sb.WriteString(" ") }
		sb.WriteString(" |")
		for _, b := range chunk {
			if b >= 32 && b < 127 { sb.WriteByte(b) } else { sb.WriteByte('.') }
		}
		sb.WriteString("|\n")
	}
	return sb.String()
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: genhex <mode> [args]
  dump   [file]         hex dump of file or stdin
  encode [text]         encode text/stdin to hex
  decode [hex]          decode hex to bytes
  rand   [n=16]         generate n random hex bytes`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 { usage() }
	mode, args := os.Args[1], os.Args[2:]
	switch mode {
	case "dump":
		var r io.Reader = os.Stdin
		if len(args) > 0 { f, err := os.Open(args[0]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }; defer f.Close(); r = f }
		data, _ := io.ReadAll(r); fmt.Print(hexDump(data))
	case "encode":
		var data []byte
		if len(args) > 0 { data = []byte(strings.Join(args, " ")) } else {
			sc := bufio.NewScanner(os.Stdin)
			for sc.Scan() { data = append(data, []byte(sc.Text()+"\n")...) }
		}
		fmt.Println(hex.EncodeToString(data))
	case "decode":
		var s string
		if len(args) > 0 { s = strings.Join(args, "") } else {
			sc := bufio.NewScanner(os.Stdin)
			for sc.Scan() { s += strings.TrimSpace(sc.Text()) }
		}
		data, err := hex.DecodeString(s); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		os.Stdout.Write(data)
	case "rand":
		n := 16
		if len(args) > 0 { n, _ = strconv.Atoi(args[0]) }
		buf := make([]byte, n); rand.Read(buf)
		fmt.Println(hex.EncodeToString(buf))
	default: usage()
	}
}
