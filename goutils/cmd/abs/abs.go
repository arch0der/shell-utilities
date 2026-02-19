// abs - print the absolute value of numbers (args or stdin, one per line)
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func process(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "abs: invalid number: %q\n", s)
		os.Exit(1)
	}
	result := math.Abs(f)
	if result == math.Trunc(result) && !strings.Contains(s, ".") {
		fmt.Printf("%.0f\n", result)
	} else {
		fmt.Printf("%g\n", result)
	}
}

func main() {
	if len(os.Args) > 1 {
		for _, a := range os.Args[1:] {
			process(a)
		}
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		process(sc.Text())
	}
}
