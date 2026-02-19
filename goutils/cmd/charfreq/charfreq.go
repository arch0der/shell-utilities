// charfreq - character frequency analysis with histogram bars
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"
)

func main() {
	freq := map[rune]int{}
	process := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			for _, ch := range sc.Text() { freq[ch]++ }
		}
	}
	if len(os.Args) == 1 {
		process(os.Stdin)
	} else {
		for _, f := range os.Args[1:] {
			fh, err := os.Open(f)
			if err != nil { fmt.Fprintln(os.Stderr, err); continue }
			process(fh); fh.Close()
		}
	}
	type kv struct{ r rune; n int }
	var pairs []kv
	total := 0
	for r, n := range freq { pairs = append(pairs, kv{r, n}); total += n }
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].n > pairs[j].n })

	maxN := 1
	if len(pairs) > 0 { maxN = pairs[0].n }
	fmt.Printf("%-8s %7s %8s  %s\n", "Char", "Count", "Percent", "Histogram")
	fmt.Println(strings.Repeat("─", 60))
	for _, p := range pairs {
		label := string(p.r)
		switch {
		case p.r == ' ': label = "SPC"
		case p.r == '\t': label = "TAB"
		case p.r == '\n': label = "LF"
		case !unicode.IsPrint(p.r): label = fmt.Sprintf("U+%04X", p.r)
		}
		pct := float64(p.n) / float64(total) * 100
		bar := strings.Repeat("█", p.n*40/maxN)
		fmt.Printf("%-8s %7d %7.2f%%  %s\n", label, p.n, pct, bar)
	}
	fmt.Printf("\nTotal: %d characters\n", total)
}
