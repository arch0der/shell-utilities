package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func init() { register("seq", runSeq) }

func runSeq() {
	args := os.Args[1:]
	separator := "\n"
	equalWidth := false
	format := ""
	nums := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-s" && i+1 < len(args):
			i++
			separator = args[i]
		case strings.HasPrefix(a, "-s"):
			separator = a[2:]
		case a == "-w" || a == "--equal-width":
			equalWidth = true
		case a == "-f" && i+1 < len(args):
			i++
			format = args[i]
		case !strings.HasPrefix(a, "-"):
			nums = append(nums, a)
		}
	}

	first, incr, last := 1.0, 1.0, 1.0
	switch len(nums) {
	case 1:
		last, _ = strconv.ParseFloat(nums[0], 64)
	case 2:
		first, _ = strconv.ParseFloat(nums[0], 64)
		last, _ = strconv.ParseFloat(nums[1], 64)
		if first > last {
			incr = -1
		}
	case 3:
		first, _ = strconv.ParseFloat(nums[0], 64)
		incr, _ = strconv.ParseFloat(nums[1], 64)
		last, _ = strconv.ParseFloat(nums[2], 64)
	default:
		fmt.Fprintln(os.Stderr, "seq: missing operand")
		os.Exit(1)
	}

	// Determine decimal places needed
	decimals := 0
	for _, n := range nums {
		if strings.Contains(n, ".") {
			parts := strings.Split(n, ".")
			if len(parts[1]) > decimals {
				decimals = len(parts[1])
			}
		}
	}

	// Calculate width for equal-width
	width := 0
	if equalWidth {
		maxVal := math.Max(math.Abs(first), math.Abs(last))
		width = len(fmt.Sprintf("%.0f", maxVal))
	}

	first2 := true
	for v := first; (incr > 0 && v <= last+1e-10) || (incr < 0 && v >= last-1e-10); v += incr {
		if !first2 {
			fmt.Print(separator)
		}
		first2 = false
		if format != "" {
			fmt.Printf(format, v)
		} else if decimals > 0 {
			fmt.Printf("%.*f", decimals, v)
		} else if equalWidth {
			fmt.Printf("%0*d", width, int64(math.Round(v)))
		} else {
			fmt.Printf("%g", v)
		}
	}
	fmt.Println()
}
