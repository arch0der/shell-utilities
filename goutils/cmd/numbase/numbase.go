// numbase - convert numbers between bases (2-36)
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: numbase <number> [from_base] [to_base]
  If from_base omitted, auto-detect from prefix (0x, 0b, 0o) or decimal.
  If to_base omitted, print all common bases.
  Bases 2-36 supported.
  examples:
    numbase 255          -> show all
    numbase 0xFF         -> auto hex, show all
    numbase 1010 2 10    -> binary to decimal
    numbase 42 10 16     -> decimal to hex`)
	os.Exit(1)
}

func autoBase(s string) int {
	switch {
	case strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X"): return 16
	case strings.HasPrefix(s, "0b") || strings.HasPrefix(s, "0B"): return 2
	case strings.HasPrefix(s, "0o") || strings.HasPrefix(s, "0O"): return 8
	default: return 10
	}
}

func stripPrefix(s string) string {
	for _, p := range []string{"0x","0X","0b","0B","0o","0O"} {
		if strings.HasPrefix(s, p) { return s[len(p):] }
	}
	return s
}

func main() {
	if len(os.Args) < 2 { usage() }
	numStr := os.Args[1]
	fromBase := autoBase(numStr)
	numStr = stripPrefix(numStr)

	if len(os.Args) >= 3 {
		n, err := strconv.Atoi(os.Args[2]); if err != nil || n < 2 || n > 36 { usage() }
		fromBase = n
	}
	val, err := strconv.ParseInt(numStr, fromBase, 64)
	if err != nil { fmt.Fprintf(os.Stderr, "numbase: cannot parse %q in base %d\n", numStr, fromBase); os.Exit(1) }

	if len(os.Args) >= 4 {
		toBase, err := strconv.Atoi(os.Args[3]); if err != nil || toBase < 2 || toBase > 36 { usage() }
		fmt.Println(strconv.FormatInt(val, toBase))
		return
	}

	fmt.Printf("Decimal  (10): %d\n", val)
	fmt.Printf("Hex      (16): 0x%X\n", val)
	fmt.Printf("Octal     (8): 0o%o\n", val)
	fmt.Printf("Binary    (2): 0b%b\n", val)
	fmt.Printf("Base36   (36): %s\n", strconv.FormatInt(val, 36))
	fmt.Printf("ASCII       : ")
	if val >= 32 && val < 127 { fmt.Printf("%c\n", val) } else { fmt.Println("(not printable ASCII)") }
}
