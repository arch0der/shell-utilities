// od - Octal/hex/decimal dump of files
// Usage: od [-t format] [-A base] [-j skip] [-N count] [file...]
// Formats: o=octal(default), x=hex, d=decimal, c=char
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	format = flag.String("t", "o", "Output format: o=octal, x=hex, d=decimal, c=char")
	addrFmt = flag.String("A", "o", "Address base: o=octal, x=hex, d=decimal, n=none")
	skip   = flag.Int64("j", 0, "Skip N bytes from start")
	count  = flag.Int64("N", -1, "Dump only N bytes")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: od [-t format] [-A addrbase] [-j skip] [-N count] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	var r io.Reader = os.Stdin
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "od:", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}

	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "od:", err)
		os.Exit(1)
	}

	if *skip > 0 && int(*skip) < len(data) {
		data = data[*skip:]
	}
	if *count >= 0 && int(*count) < len(data) {
		data = data[:*count]
	}

	lineWidth := 16
	for offset := 0; offset < len(data); offset += lineWidth {
		end := offset + lineWidth
		if end > len(data) {
			end = len(data)
		}
		chunk := data[offset:end]

		// Print address
		switch *addrFmt {
		case "x":
			fmt.Printf("%06x", offset)
		case "d":
			fmt.Printf("%07d", offset)
		case "n":
			// no address
		default: // "o"
			fmt.Printf("%07o", offset)
		}

		// Print values
		for _, b := range chunk {
			switch *format {
			case "x":
				fmt.Printf(" %02x", b)
			case "d":
				fmt.Printf(" %3d", b)
			case "c":
				if b >= 0x20 && b < 0x7f {
					fmt.Printf("   %c", b)
				} else {
					fmt.Printf(" %03o", b)
				}
			default: // "o"
				fmt.Printf(" %03o", b)
			}
		}
		fmt.Println()
	}

	// Final address
	switch *addrFmt {
	case "x":
		fmt.Printf("%06x\n", len(data))
	case "d":
		fmt.Printf("%07d\n", len(data))
	case "n":
	default:
		fmt.Printf("%07o\n", len(data))
	}
}
