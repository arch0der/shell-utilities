package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func init() { register("od", runOd) }

func runOd() {
	args := os.Args[1:]
	format := "o2"  // default: octal shorts
	addrFmt := "o"
	files := []string{}
	skipBytes := int64(0)

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-A" && i+1 < len(args):
			i++
			addrFmt = args[i]
		case a == "-t" && i+1 < len(args):
			i++
			format = args[i]
		case strings.HasPrefix(a, "-t"):
			format = a[2:]
		case a == "-x":
			format = "x2"
		case a == "-o":
			format = "o2"
		case a == "-d":
			format = "u2"
		case a == "-c":
			format = "c"
		case a == "-b":
			format = "o1"
		case a == "-j" && i+1 < len(args):
			i++
			fmt.Sscan(args[i], &skipBytes)
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}

	var r io.Reader = os.Stdin
	if len(files) > 0 {
		fh, err := os.Open(files[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "od: %v\n", err)
			os.Exit(1)
		}
		defer fh.Close()
		if skipBytes > 0 {
			fh.Seek(skipBytes, io.SeekStart)
		}
		r = fh
	}

	data, _ := io.ReadAll(r)

	printAddr := func(addr int) {
		switch addrFmt {
		case "o":
			fmt.Printf("%07o", addr)
		case "x":
			fmt.Printf("%06x", addr)
		case "d":
			fmt.Printf("%07d", addr)
		case "n":
			// no address
		}
	}

	width := 16
	for i := 0; i < len(data); i += width {
		end := i + width
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]
		printAddr(i)

		switch format {
		case "c":
			for _, b := range chunk {
				switch b {
				case '\n':
					fmt.Print("  \\n")
				case '\t':
					fmt.Print("  \\t")
				case '\r':
					fmt.Print("  \\r")
				case 0:
					fmt.Print(" \\00")
				default:
					if b >= 32 && b < 127 {
						fmt.Printf("   %c", b)
					} else {
						fmt.Printf(" %03o", b)
					}
				}
			}
		case "o1":
			for _, b := range chunk {
				fmt.Printf(" %03o", b)
			}
		case "x1":
			for _, b := range chunk {
				fmt.Printf(" %02x", b)
			}
		case "x2":
			for j := 0; j < len(chunk); j += 2 {
				if j+1 < len(chunk) {
					fmt.Printf(" %04x", uint16(chunk[j])|uint16(chunk[j+1])<<8)
				} else {
					fmt.Printf(" %04x", uint16(chunk[j]))
				}
			}
		case "o2":
			for j := 0; j < len(chunk); j += 2 {
				if j+1 < len(chunk) {
					fmt.Printf(" %06o", uint16(chunk[j])|uint16(chunk[j+1])<<8)
				} else {
					fmt.Printf(" %06o", uint16(chunk[j]))
				}
			}
		case "u2":
			for j := 0; j < len(chunk); j += 2 {
				if j+1 < len(chunk) {
					fmt.Printf(" %5d", uint16(chunk[j])|uint16(chunk[j+1])<<8)
				} else {
					fmt.Printf(" %5d", uint16(chunk[j]))
				}
			}
		}
		fmt.Println()
	}
	printAddr(len(data))
	fmt.Println()
}
