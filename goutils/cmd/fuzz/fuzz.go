// fuzz - generate random fuzzing inputs (strings, numbers, bytes)
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	printable = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !@#$%^&*()_+-=[]{}|;':\",./<>?"
	unicode_specials = []rune{0, 1, 127, 0xFF, 0x100, 0xFFFD, 0x1F4A9, 0x202E, 0xFEFF}
)

func randStr(rng *rand.Rand, minLen, maxLen int) string {
	n := minLen + rng.Intn(maxLen-minLen+1)
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(printable[rng.Intn(len(printable))])
	}
	return b.String()
}

func randBytes(rng *rand.Rand, n int) string {
	buf := make([]byte, n)
	rng.Read(buf)
	var b strings.Builder
	for _, c := range buf { fmt.Fprintf(&b, "\\x%02x", c) }
	return b.String()
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: fuzz <type> [count] [options]
  types: string | int | float | bytes | unicode | boundary
  string count [minlen] [maxlen]
  int    count [min] [max]
  float  count [min] [max]
  bytes  count [nbytes]
  unicode count
  boundary  (common edge-case strings)`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 { usage() }
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	typ := os.Args[1]
	count := 10
	if len(os.Args) > 2 { count, _ = strconv.Atoi(os.Args[2]) }
	if count < 1 { count = 1 }

	switch typ {
	case "string":
		min, max := 1, 64
		if len(os.Args) > 3 { min, _ = strconv.Atoi(os.Args[3]) }
		if len(os.Args) > 4 { max, _ = strconv.Atoi(os.Args[4]) }
		for i := 0; i < count; i++ { fmt.Println(randStr(rng, min, max)) }
	case "int":
		lo, hi := int64(-1000000), int64(1000000)
		if len(os.Args) > 3 { lo, _ = strconv.ParseInt(os.Args[3], 10, 64) }
		if len(os.Args) > 4 { hi, _ = strconv.ParseInt(os.Args[4], 10, 64) }
		for i := 0; i < count; i++ { fmt.Println(lo + rng.Int63n(hi-lo+1)) }
	case "float":
		lo, hi := -1000.0, 1000.0
		if len(os.Args) > 3 { lo, _ = strconv.ParseFloat(os.Args[3], 64) }
		if len(os.Args) > 4 { hi, _ = strconv.ParseFloat(os.Args[4], 64) }
		for i := 0; i < count; i++ { fmt.Printf("%g\n", lo+rng.Float64()*(hi-lo)) }
	case "bytes":
		nb := 16
		if len(os.Args) > 3 { nb, _ = strconv.Atoi(os.Args[3]) }
		for i := 0; i < count; i++ { fmt.Println(randBytes(rng, nb)) }
	case "unicode":
		for i := 0; i < count; i++ { fmt.Printf("%c (U+%04X)\n", unicode_specials[i%len(unicode_specials)], unicode_specials[i%len(unicode_specials)]) }
	case "boundary":
		boundaries := []string{
			"", " ", "  ", "\t", "\n", "\r\n", "null", "NULL", "nil", "None",
			"0", "-1", "1", "999999999", "-999999999", "2147483647", "-2147483648",
			"<script>alert(1)</script>", "' OR '1'='1", `\"`, "\\", "/",
			strings.Repeat("A", 256), strings.Repeat("A", 65536),
			"café", "こんにちは", "🎉", "\x00", "\xFF",
		}
		for _, b := range boundaries { fmt.Printf("%q\n", b) }
	default:
		fmt.Fprintf(os.Stderr, "fuzz: unknown type %q\n", typ); usage()
	}
}
