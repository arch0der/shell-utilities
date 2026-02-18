// diff - Compare two files line by line
// Usage: diff [-u] [-i] <file1> <file2>
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	unified    = flag.Bool("u", false, "Unified diff format")
	ignoreCase = flag.Bool("i", false, "Ignore case differences")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: diff [-u] [-i] <file1> <file2>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	a, err := readLines(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "diff:", err)
		os.Exit(1)
	}
	b, err := readLines(flag.Arg(1))
	if err != nil {
		fmt.Fprintln(os.Stderr, "diff:", err)
		os.Exit(1)
	}

	if *ignoreCase {
		for i := range a {
			a[i] = strings.ToLower(a[i])
		}
		for i := range b {
			b[i] = strings.ToLower(b[i])
		}
	}

	hunks := lcs(a, b)
	if len(hunks) == 0 {
		os.Exit(0)
	}

	if *unified {
		fmt.Printf("--- %s\n", flag.Arg(0))
		fmt.Printf("+++ %s\n", flag.Arg(1))
	}

	for _, h := range hunks {
		if *unified {
			fmt.Printf("@@ -%d,%d +%d,%d @@\n", h.a1+1, h.a2-h.a1, h.b1+1, h.b2-h.b1)
			for i := h.a1; i < h.a2; i++ {
				fmt.Println("-" + a[i])
			}
			for i := h.b1; i < h.b2; i++ {
				fmt.Println("+" + b[i])
			}
		} else {
			if h.a1 == h.a2 {
				fmt.Printf("%da%d,%d\n", h.a1, h.b1+1, h.b2)
			} else if h.b1 == h.b2 {
				fmt.Printf("%d,%dd%d\n", h.a1+1, h.a2, h.b1)
			} else {
				fmt.Printf("%d,%dc%d,%d\n", h.a1+1, h.a2, h.b1+1, h.b2)
			}
			for i := h.a1; i < h.a2; i++ {
				fmt.Println("< " + a[i])
			}
			if h.a1 != h.a2 && h.b1 != h.b2 {
				fmt.Println("---")
			}
			for i := h.b1; i < h.b2; i++ {
				fmt.Println("> " + b[i])
			}
		}
	}
	os.Exit(1)
}

type hunk struct{ a1, a2, b1, b2 int }

// Simple O(n*m) diff using edit distance
func lcs(a, b []string) []hunk {
	m, n := len(a), len(b)
	// Build edit matrix
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else {
				if dp[i+1][j] > dp[i][j+1] {
					dp[i][j] = dp[i+1][j]
				} else {
					dp[i][j] = dp[i][j+1]
				}
			}
		}
	}

	// Trace back to find hunks
	var hunks []hunk
	i, j := 0, 0
	var cur *hunk
	for i < m || j < n {
		if i < m && j < n && a[i] == b[j] {
			if cur != nil {
				hunks = append(hunks, *cur)
				cur = nil
			}
			i++
			j++
		} else {
			if cur == nil {
				cur = &hunk{a1: i, a2: i, b1: j, b2: j}
			}
			if j >= n || (i < m && dp[i+1][j] >= dp[i][j+1]) {
				cur.a2 = i + 1
				i++
			} else {
				cur.b2 = j + 1
				j++
			}
		}
	}
	if cur != nil {
		hunks = append(hunks, *cur)
	}
	return hunks
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
