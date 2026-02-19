package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func main() {
	args := os.Args[1:]
	files := []string{}
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}

	var words []string
	var lines []string

	readFile := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			line := sc.Text()
			lines = append(lines, line)
			for _, w := range strings.Fields(line) {
				words = append(words, strings.ToLower(w))
			}
		}
	}

	if len(files) == 0 {
		readFile(os.Stdin)
	}
	for _, f := range files {
		fh, _ := os.Open(f)
		readFile(fh)
		fh.Close()
	}

	// Build permuted index
	type entry struct {
		before, word, after string
	}
	var entries []entry
	for _, line := range lines {
		words2 := strings.Fields(line)
		for i, w := range words2 {
			before := strings.Join(words2[:i], " ")
			after := strings.Join(words2[i+1:], " ")
			entries = append(entries, entry{before, w, after})
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].word) < strings.ToLower(entries[j].word)
	})
	for _, e := range entries {
		fmt.Printf("%-20s %s /%s/\n", e.before, e.word, e.after)
	}
	_ = words
}
