// histogram - plot a text histogram from numeric input (stdin, one number per line)
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	bins := 10
	width := 60
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-b": i++; bins, _ = strconv.Atoi(args[i])
		case "-w": i++; width, _ = strconv.Atoi(args[i])
		}
	}

	var vals []float64
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s == "" { continue }
		f, err := strconv.ParseFloat(s, 64)
		if err != nil { continue }
		vals = append(vals, f)
	}
	if len(vals) == 0 { fmt.Fprintln(os.Stderr, "histogram: no data"); os.Exit(1) }

	mn, mx := vals[0], vals[0]
	sum := 0.0
	for _, v := range vals {
		if v < mn { mn = v }
		if v > mx { mx = v }
		sum += v
	}
	mean := sum / float64(len(vals))
	rng := mx - mn
	if rng == 0 { rng = 1 }
	binSize := rng / float64(bins)

	counts := make([]int, bins)
	for _, v := range vals {
		b := int((v - mn) / binSize)
		if b >= bins { b = bins - 1 }
		counts[b]++
	}
	maxCount := 0
	for _, c := range counts { if c > maxCount { maxCount = c } }

	fmt.Printf("n=%d  min=%.4g  max=%.4g  mean=%.4g  bins=%d\n\n", len(vals), mn, mx, mean, bins)
	for i, c := range counts {
		lo := mn + float64(i)*binSize
		hi := lo + binSize
		barLen := 0
		if maxCount > 0 { barLen = c * width / maxCount }
		bar := strings.Repeat("█", barLen)
		fmt.Printf("[%9.4g, %9.4g) │%-*s %d\n", lo, hi, width, bar, c)
	}
	fmt.Printf("\nStdDev: %.4g\n", stddev(vals, mean))
}

func stddev(vals []float64, mean float64) float64 {
	sum := 0.0
	for _, v := range vals { d := v - mean; sum += d * d }
	return math.Sqrt(sum / float64(len(vals)))
}
