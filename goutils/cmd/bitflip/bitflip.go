// bitflip - flip bits in an integer; all bits or specific positions (0=LSB)
package main

import (
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func binStr(n uint64, width int) string {
	s := strconv.FormatUint(n, 2)
	for len(s) < width { s = "0" + s }
	// group by 4
	var b strings.Builder
	for i, ch := range s {
		if i > 0 && (len(s)-i)%4 == 0 { b.WriteRune('_') }
		b.WriteRune(ch)
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: bitflip <number> [bit_pos ...]")
		os.Exit(1)
	}
	n, err := strconv.ParseUint(strings.TrimPrefix(os.Args[1], "0x"), 16, 64)
	if err != nil {
		n64, err2 := strconv.ParseInt(os.Args[1], 0, 64)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "bitflip: invalid number: %s\n", os.Args[1])
			os.Exit(1)
		}
		n = uint64(n64)
	} else if !strings.HasPrefix(os.Args[1], "0x") {
		n64, _ := strconv.ParseInt(os.Args[1], 0, 64)
		n = uint64(n64)
	}

	w := bits.Len64(n)
	if w < 8 { w = 8 }

	var result uint64
	if len(os.Args) == 2 {
		result = ^n
	} else {
		result = n
		for _, arg := range os.Args[2:] {
			bit, err := strconv.Atoi(arg)
			if err != nil || bit < 0 || bit > 63 {
				fmt.Fprintf(os.Stderr, "bitflip: invalid bit position: %s\n", arg)
				os.Exit(1)
			}
			result ^= (1 << uint(bit))
		}
	}
	fmt.Printf("Original : 0b%s  (dec:%d  hex:0x%X)\n", binStr(n, w), n, n)
	fmt.Printf("Result   : 0b%s  (dec:%d  hex:0x%X)\n", binStr(result, w), result, result)
	fmt.Printf("Changed  : 0b%s\n", binStr(n^result, w))
}
