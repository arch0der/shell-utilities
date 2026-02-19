// matrix - matrix math operations on stdin data
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func readMatrix(r *os.File) [][]float64 {
	var m [][]float64
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") { continue }
		parts := strings.Fields(line)
		row := make([]float64, len(parts))
		for i, p := range parts {
			v, err := strconv.ParseFloat(p, 64)
			if err != nil { fmt.Fprintf(os.Stderr, "matrix: bad value %q\n", p); os.Exit(1) }
			row[i] = v
		}
		m = append(m, row)
	}
	return m
}

func printMatrix(m [][]float64) {
	for _, row := range m {
		parts := make([]string, len(row))
		for i, v := range row { parts[i] = fmt.Sprintf("%10.4g", v) }
		fmt.Println(strings.Join(parts, "  "))
	}
}

func transpose(m [][]float64) [][]float64 {
	if len(m) == 0 { return nil }
	rows, cols := len(m), len(m[0])
	t := make([][]float64, cols)
	for i := range t { t[i] = make([]float64, rows) }
	for i, row := range m { for j, v := range row { t[j][i] = v } }
	return t
}

func multiply(a, b [][]float64) [][]float64 {
	ra, ca := len(a), len(a[0])
	rb, cb := len(b), len(b[0])
	if ca != rb { fmt.Fprintln(os.Stderr, "matrix: dimension mismatch for multiply"); os.Exit(1) }
	c := make([][]float64, ra)
	for i := range c {
		c[i] = make([]float64, cb)
		for j := 0; j < cb; j++ {
			for k := 0; k < ca; k++ { c[i][j] += a[i][k] * b[k][j] }
		}
	}
	_ = rb; return c
}

func stats(m [][]float64) {
	var vals []float64
	for _, row := range m { vals = append(vals, row...) }
	sum, mn, mx := 0.0, vals[0], vals[0]
	for _, v := range vals { sum += v; if v < mn { mn = v }; if v > mx { mx = v } }
	mean := sum / float64(len(vals))
	variance := 0.0
	for _, v := range vals { d := v - mean; variance += d * d }
	fmt.Printf("Rows: %d  Cols: %d  Elements: %d\n", len(m), len(m[0]), len(vals))
	fmt.Printf("Min: %g  Max: %g  Sum: %g  Mean: %g  StdDev: %g\n",
		mn, mx, sum, mean, math.Sqrt(variance/float64(len(vals))))
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: matrix <op> [file1] [file2]")
	fmt.Fprintln(os.Stderr, "  ops: transpose | multiply | add | sub | stats")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 { usage() }
	op := os.Args[1]
	var m1, m2 [][]float64

	getFile := func(idx int) *os.File {
		if idx < len(os.Args) {
			f, err := os.Open(os.Args[idx])
			if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
			return f
		}
		return os.Stdin
	}

	switch op {
	case "transpose":
		m1 = readMatrix(getFile(2))
		printMatrix(transpose(m1))
	case "multiply":
		m1 = readMatrix(getFile(2)); m2 = readMatrix(getFile(3))
		printMatrix(multiply(m1, m2))
	case "add", "sub":
		m1 = readMatrix(getFile(2)); m2 = readMatrix(getFile(3))
		if len(m1) != len(m2) || len(m1[0]) != len(m2[0]) { fmt.Fprintln(os.Stderr, "matrix: dimension mismatch"); os.Exit(1) }
		for i, row := range m1 { for j := range row {
			if op == "add" { m1[i][j] += m2[i][j] } else { m1[i][j] -= m2[i][j] }
		}}
		printMatrix(m1)
	case "stats":
		m1 = readMatrix(getFile(2)); stats(m1)
	default: usage()
	}
}
