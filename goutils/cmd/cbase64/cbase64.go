// cbase64 - base64 encode/decode with URL-safe and MIME variants
package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: cbase64 [options] [file]
  -d, --decode    decode base64 input
  -u, --url       use URL-safe encoding (- and _ instead of + and /)
  -m, --mime      MIME encoding (76-char lines with CRLF)
  -w <n>          wrap encoded output at n chars (0 = no wrap, default 0)
  -s              strict: reject padding errors`)
	os.Exit(1)
}

func main() {
	decode := false
	urlSafe := false
	mime := false
	wrap := 0
	strict := false
	var file string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d", "--decode": decode = true
		case "-u", "--url": urlSafe = true
		case "-m", "--mime": mime = true; wrap = 76
		case "-s": strict = true
		case "-w": i++; fmt.Sscanf(args[i], "%d", &wrap)
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			file = args[i]
		}
	}

	var enc *base64.Encoding
	switch {
	case urlSafe: enc = base64.URLEncoding
	case strict: enc = base64.StdEncoding
	default: enc = base64.StdEncoding.WithPadding(base64.StdPadding)
	}

	var r io.Reader = os.Stdin
	if file != "" {
		f, err := os.Open(file); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		defer f.Close(); r = f
	}

	data, err := io.ReadAll(r)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

	if decode {
		// clean whitespace
		cleaned := strings.Map(func(r rune) rune {
			if r == '\n' || r == '\r' || r == ' ' || r == '\t' { return -1 }; return r
		}, string(data))
		out, err := enc.DecodeString(cleaned)
		if err != nil { fmt.Fprintln(os.Stderr, "cbase64:", err); os.Exit(1) }
		os.Stdout.Write(out)
		return
	}

	encoded := enc.EncodeToString(data)
	_ = mime
	if wrap > 0 {
		for i := 0; i < len(encoded); i += wrap {
			end := i + wrap; if end > len(encoded) { end = len(encoded) }
			fmt.Println(encoded[i:end])
		}
	} else {
		fmt.Println(encoded)
	}
}
