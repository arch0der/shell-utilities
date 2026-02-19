// b64url - URL-safe Base64 encode/decode (RFC 4648).
//
// Usage:
//
//	b64url [OPTIONS] [FILE...]
//	echo "hello world" | b64url
//	echo "aGVsbG8gd29ybGQ" | b64url -d
//
// Options:
//
//	-d        Decode
//	-s        Standard base64 (not URL-safe)
//	-n        No padding (omit = characters)
//	-w N      Wrap encoded output at N chars (default: 0 = no wrap)
//
// Examples:
//
//	echo "hello" | b64url           # aGVsbG8=  (URL-safe)
//	echo "aGVsbG8" | b64url -d      # hello
//	cat file.bin | b64url -n        # no padding
//	echo "hello" | b64url -s        # standard base64
package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	decode   = flag.Bool("d", false, "decode")
	standard = flag.Bool("s", false, "standard base64")
	noPad    = flag.Bool("n", false, "no padding")
	wrap     = flag.Int("w", 0, "wrap width")
)

func encoder() *base64.Encoding {
	var enc *base64.Encoding
	if *standard {
		enc = base64.StdEncoding
	} else {
		enc = base64.URLEncoding
	}
	if *noPad {
		enc = enc.WithPadding(base64.NoPadding)
	}
	return enc
}

func wrapStr(s string, width int) string {
	if width <= 0 {
		return s
	}
	var sb strings.Builder
	for i := 0; i < len(s); i += width {
		end := i + width
		if end > len(s) {
			end = len(s)
		}
		sb.WriteString(s[i:end])
		if end < len(s) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func main() {
	flag.Parse()
	files := flag.Args()

	var readers []io.Reader
	if len(files) == 0 {
		readers = []io.Reader{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "b64url: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}

	enc := encoder()
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for _, r := range readers {
		data, err := io.ReadAll(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "b64url: %v\n", err)
			os.Exit(1)
		}

		if *decode {
			// Strip whitespace
			clean := strings.ReplaceAll(string(data), "\n", "")
			clean = strings.ReplaceAll(clean, "\r", "")
			clean = strings.TrimSpace(clean)
			decoded, err := enc.DecodeString(clean)
			if err != nil {
				// Try with padding
				padded := clean
				switch len(clean) % 4 {
				case 2:
					padded += "=="
				case 3:
					padded += "="
				}
				decoded, err = enc.DecodeString(padded)
				if err != nil {
					fmt.Fprintf(os.Stderr, "b64url: decode error: %v\n", err)
					os.Exit(1)
				}
			}
			w.Write(decoded)
		} else {
			encoded := enc.EncodeToString(data)
			if *wrap > 0 {
				encoded = wrapStr(encoded, *wrap)
			}
			fmt.Fprintln(w, encoded)
		}
	}
}
