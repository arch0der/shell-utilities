// seq - Print a sequence of numbers
// Usage: seq [-w] [-s separator] [-f format] [first [step]] last
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
)

var (
	pad    = flag.Bool("w", false, "Equalize width by padding with leading zeros")
	sep    = flag.String("s", "\n", "Separator between numbers")
	format = flag.String("f", "", "Printf format (e.g. %05.2f)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: seq [-w] [-s sep] [-f fmt] [first [step]] last")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var first, step, last float64
	switch len(args) {
	case 1:
		first = 1
		step = 1
		last, _ = strconv.ParseFloat(args[0], 64)
	case 2:
		first, _ = strconv.ParseFloat(args[0], 64)
		step = 1
		last, _ = strconv.ParseFloat(args[1], 64)
		if first > last {
			step = -1
		}
	default:
		first, _ = strconv.ParseFloat(args[0], 64)
		step, _ = strconv.ParseFloat(args[1], 64)
		last, _ = strconv.ParseFloat(args[2], 64)
	}

	if step == 0 {
		fmt.Fprintln(os.Stderr, "seq: zero increment")
		os.Exit(1)
	}

	// Determine format
	fmtStr := *format
	if fmtStr == "" {
		// Determine decimal places
		decFirst := decimals(first)
		decStep := decimals(step)
		dec := decFirst
		if decStep > dec {
			dec = decStep
		}
		if dec == 0 {
			fmtStr = "%g"
			if *pad {
				// Figure out max width
				maxVal := last
				if first > last {
					maxVal = first
				}
				width := len(fmt.Sprintf("%g", maxVal))
				fmtStr = fmt.Sprintf("%%0%dg", width)
			}
		} else {
			fmtStr = fmt.Sprintf("%%.%df", dec)
		}
	}

	first_result := true
	for n := first; (step > 0 && n <= last+1e-10) || (step < 0 && n >= last-1e-10); n += step {
		if !first_result {
			fmt.Print(*sep)
		}
		first_result = false
		fmt.Printf(fmtStr, n)
	}
	fmt.Println()
}

func decimals(f float64) int {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	for i, c := range s {
		if c == '.' {
			return len(s) - i - 1
		}
	}
	_ = math.Pi // avoid unused import
	return 0
}
