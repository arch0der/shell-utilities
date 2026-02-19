// uuid - Generate UUIDs (v4 random)
// Usage: uuid [-n count] [-upper]
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	count = flag.Int("n", 1, "Number of UUIDs to generate")
	upper = flag.Bool("upper", false, "Output in uppercase")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: uuid [-n count] [-upper]")
		flag.PrintDefaults()
	}
	flag.Parse()

	for i := 0; i < *count; i++ {
		u, err := newUUID()
		if err != nil {
			fmt.Fprintln(os.Stderr, "uuid:", err)
			os.Exit(1)
		}
		if *upper {
			fmt.Println(strings.ToUpper(u))
		} else {
			fmt.Println(u)
		}
	}
}

func newUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Set version 4
	b[6] = (b[6] & 0x0f) | 0x40
	// Set variant bits (10xx)
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
