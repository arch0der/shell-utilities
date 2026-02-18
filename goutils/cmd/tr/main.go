// tr - Translate or delete characters
// Usage: tr [-d] [-s] <set1> [set2]
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	deleteMode  = flag.Bool("d", false, "Delete characters in set1")
	squeezeMode = flag.Bool("s", false, "Squeeze repeated characters in set1")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: tr [-d] [-s] <set1> [set2]")
		fmt.Fprintln(os.Stderr, "Escape sequences: \\n \\t \\r; ranges: a-z; classes: [:upper:] [:lower:] [:digit:]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	set1 := expandSet(flag.Arg(0))
	set2 := ""
	if flag.NArg() >= 2 {
		set2 = expandSet(flag.Arg(1))
	}

	// Build translation table
	table := map[rune]rune{}
	if !*deleteMode && set2 != "" {
		r2 := []rune(set2)
		for i, c := range []rune(set1) {
			if i < len(r2) {
				table[c] = r2[i]
			} else {
				table[c] = r2[len(r2)-1]
			}
		}
	}

	deleteSet := map[rune]bool{}
	if *deleteMode {
		for _, c := range []rune(set1) {
			deleteSet[c] = true
		}
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tr:", err)
		os.Exit(1)
	}

	var out strings.Builder
	var prev rune = -1
	for _, c := range string(data) {
		if *deleteMode && deleteSet[c] {
			continue
		}
		if t, ok := table[c]; ok {
			c = t
		}
		if *squeezeMode && c == prev {
			_, inSet1 := map[rune]bool{}[c]
			for _, r := range []rune(set1) {
				if r == c {
					inSet1 = true
					break
				}
			}
			if inSet1 {
				continue
			}
		}
		out.WriteRune(c)
		prev = c
	}
	fmt.Print(out.String())
}

func expandSet(s string) string {
	// Handle escape sequences
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\t`, "\t")
	s = strings.ReplaceAll(s, `\r`, "\r")
	s = strings.ReplaceAll(s, `\\`, "\\")

	// Handle POSIX classes
	s = strings.ReplaceAll(s, "[:upper:]", "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s = strings.ReplaceAll(s, "[:lower:]", "abcdefghijklmnopqrstuvwxyz")
	s = strings.ReplaceAll(s, "[:digit:]", "0123456789")
	s = strings.ReplaceAll(s, "[:alpha:]", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	s = strings.ReplaceAll(s, "[:space:]", " \t\n\r")

	// Handle ranges like a-z
	var out strings.Builder
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if i+2 < len(runes) && runes[i+1] == '-' {
			from, to := runes[i], runes[i+2]
			for c := from; c <= to; c++ {
				out.WriteRune(c)
			}
			i += 2
		} else {
			out.WriteRune(runes[i])
		}
	}
	return out.String()
}
