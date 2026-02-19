// csvstat - compute summary statistics for each numeric column in a CSV
package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type colStats struct {
	name    string
	vals    []float64
	nonNum  int
	empty   int
}

func (c *colStats) compute() {
	if len(c.vals) == 0 { return }
	sort.Float64s(c.vals)
	n := float64(len(c.vals))
	sum := 0.0
	for _, v := range c.vals { sum += v }
	mean := sum / n
	variance := 0.0
	for _, v := range c.vals { d := v - mean; variance += d * d }
	stddev := math.Sqrt(variance / n)
	median := c.vals[len(c.vals)/2]
	if len(c.vals)%2 == 0 { median = (c.vals[len(c.vals)/2-1] + c.vals[len(c.vals)/2]) / 2 }

	fmt.Printf("%-20s  n=%-6d  min=%-12g  max=%-12g  mean=%-12g  median=%-12g  stddev=%g\n",
		c.name, len(c.vals), c.vals[0], c.vals[len(c.vals)-1], mean, median, stddev)
}

func main() {
	var file string
	if len(os.Args) > 1 { file = os.Args[1] }

	var r *os.File
	if file != "" {
		var err error; r, err = os.Open(file)
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		defer r.Close()
	} else { r = os.Stdin }

	cr := csv.NewReader(r); cr.LazyQuotes = true; cr.TrimLeadingSpace = true
	recs, err := cr.ReadAll()
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	if len(recs) < 2 { fmt.Println("Not enough data"); return }

	headers := recs[0]
	cols := make([]*colStats, len(headers))
	for i, h := range headers { cols[i] = &colStats{name: h} }

	for _, row := range recs[1:] {
		for i, cell := range row {
			if i >= len(cols) { continue }
			cell = strings.TrimSpace(cell)
			if cell == "" { cols[i].empty++; continue }
			v, err := strconv.ParseFloat(cell, 64)
			if err != nil { cols[i].nonNum++; continue }
			cols[i].vals = append(cols[i].vals, v)
		}
	}

	fmt.Printf("File: %s  Rows: %d\n\n", func() string { if file == "" { return "<stdin>" }; return file }(), len(recs)-1)
	fmt.Printf("%-20s  %-8s  %-14s %-14s %-14s %-14s %s\n",
		"Column", "Count", "Min", "Max", "Mean", "Median", "StdDev")
	fmt.Println(strings.Repeat("─", 100))
	for _, c := range cols { c.compute() }
}
