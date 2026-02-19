// linediff - show side-by-side diff of two files with colour highlighting
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
	dim   = "\033[2m"
)

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil { return nil, err }
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() { lines = append(lines, sc.Text()) }
	return lines, nil
}

func lcs(a, b []string) [][]int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp { dp[i] = make([]int, n+1) }
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] { dp[i][j] = dp[i-1][j-1] + 1 } else if dp[i-1][j] > dp[i][j-1] { dp[i][j] = dp[i-1][j] } else { dp[i][j] = dp[i][j-1] }
		}
	}
	return dp
}

type change struct{ op rune; line string }

func diff(a, b []string) []change {
	dp := lcs(a, b)
	var changes []change
	i, j := len(a), len(b)
	for i > 0 || j > 0 {
		switch {
		case i > 0 && j > 0 && a[i-1] == b[j-1]:
			changes = append([]change{{'=', a[i-1]}}, changes...); i--; j--
		case j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]):
			changes = append([]change{{'+', b[j-1]}}, changes...); j--
		default:
			changes = append([]change{{'-', a[i-1]}}, changes...); i--
		}
	}
	return changes
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: linediff <file1> <file2>"); os.Exit(1)
	}
	a, err := readLines(os.Args[1])
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	b, err2 := readLines(os.Args[2])
	if err2 != nil { fmt.Fprintln(os.Stderr, err2); os.Exit(1) }

	changes := diff(a, b)
	w := 60
	fmt.Printf("%-*s │ %s\n", w, os.Args[1], os.Args[2])
	fmt.Println(strings.Repeat("─", w) + "─┼─" + strings.Repeat("─", w))
	for _, ch := range changes {
		line := ch.line
		if len(line) > w-3 { line = line[:w-3] + "..." }
		switch ch.op {
		case '=': fmt.Printf("%-*s │ %s\n", w, line, line)
		case '-': fmt.Printf("%s%-*s%s │ %s(removed)%s\n", red, w, line, reset, dim, reset)
		case '+': fmt.Printf("%-*s │ %s%s%s\n", w, dim+"(added)"+reset, green, line, reset)
		}
	}
}
