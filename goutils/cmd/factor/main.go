// factor - Print prime factorization of numbers
// Usage: factor [number...]
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			n, err := strconv.ParseUint(arg, 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "factor: invalid number: %s\n", arg)
				continue
			}
			printFactors(n)
		}
		return
	}

	// Read from stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		for _, f := range strings.Fields(sc.Text()) {
			n, err := strconv.ParseUint(f, 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "factor: invalid number: %s\n", f)
				continue
			}
			printFactors(n)
		}
	}
}

func printFactors(n uint64) {
	fmt.Printf("%d:", n)
	if n <= 1 {
		fmt.Printf(" %d", n)
		fmt.Println()
		return
	}
	orig := n
	_ = orig
	for p := uint64(2); p*p <= n; p++ {
		for n%p == 0 {
			fmt.Printf(" %d", p)
			n /= p
		}
	}
	if n > 1 {
		fmt.Printf(" %d", n)
	}
	fmt.Println()
}
