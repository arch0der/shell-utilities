// bcalc - bitwise calculator: AND, OR, XOR, NOT, shifts on integers
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseNum(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "0b") || strings.HasPrefix(s, "0B") { return strconv.ParseInt(s[2:], 2, 64) }
	return strconv.ParseInt(s, 0, 64)
}

func fmtNum(n int64) {
	u := uint64(n)
	fmt.Printf("dec  : %d\n", n)
	fmt.Printf("hex  : 0x%X\n", u)
	fmt.Printf("oct  : 0o%o\n", u)
	fmt.Printf("bin  : 0b%b\n", u)
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: bcalc <a> <op> <b>   |   bcalc not <a>
  ops: and | or | xor | nand | nor | shl | shr | rotl | rotr
  values can be: 0b1010 (binary), 0xFF (hex), 0o77 (octal), 42 (decimal)`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 { usage() }
	if strings.ToLower(os.Args[1]) == "not" {
		a, err := parseNum(os.Args[2]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		fmtNum(^a); return
	}
	if len(os.Args) < 4 { usage() }
	a, err := parseNum(os.Args[1]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	op := strings.ToLower(os.Args[2])
	b, err := parseNum(os.Args[3]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

	var result int64
	switch op {
	case "and": result = a & b
	case "or":  result = a | b
	case "xor": result = a ^ b
	case "nand": result = ^(a & b)
	case "nor":  result = ^(a | b)
	case "shl", "<<": result = a << uint(b)
	case "shr", ">>": result = a >> uint(b)
	case "rotl": result = int64(rotl64(uint64(a), uint(b)))
	case "rotr": result = int64(rotr64(uint64(a), uint(b)))
	default: fmt.Fprintf(os.Stderr, "bcalc: unknown op %q\n", op); usage()
	}
	fmt.Printf("%-6s: %d %s %d\n", "expr", a, op, b)
	fmtNum(result)
}

func rotl64(x uint64, n uint) uint64 { n &= 63; return (x<<n)|(x>>(64-n)) }
func rotr64(x uint64, n uint) uint64 { n &= 63; return (x>>n)|(x<<(64-n)) }
