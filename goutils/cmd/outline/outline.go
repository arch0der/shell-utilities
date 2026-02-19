// outline - extract heading structure from Markdown or plain text (outline view)
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var headingRe = regexp.MustCompile(`^(#{1,6})\s+(.+)`)
var numberedRe = regexp.MustCompile(`^(\s*)(\d+[.)]\s+|\*\s+|-\s+)(.+)`)

func main() {
	flat := false
	maxLevel := 6
	var files []string
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-f", "--flat": flat = true
		case "--h1": maxLevel = 1
		case "--h2": maxLevel = 2
		case "--h3": maxLevel = 3
		default:
			if !strings.HasPrefix(arg, "-") { files = append(files, arg) }
		}
	}

	process := func(r *os.File) {
		sc := bufio.NewScanner(r)
		counters := make([]int, 7)
		for sc.Scan() {
			line := sc.Text()
			m := headingRe.FindStringSubmatch(line)
			if m == nil { continue }
			level := len(m[1])
			if level > maxLevel { continue }
			text := strings.TrimSpace(m[2])
			// strip markdown formatting from heading
			text = regexp.MustCompile(`\*\*?|__?|~~|` + "`").ReplaceAllString(text, "")
			text = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`).ReplaceAllString(text, "$1")

			if flat {
				fmt.Printf("%s %s\n", strings.Repeat("#", level), text)
				continue
			}

			// Update counters
			counters[level]++
			for i := level + 1; i <= 6; i++ { counters[i] = 0 }

			// Build number
			var parts []string
			for i := 1; i <= level; i++ { if counters[i] > 0 { parts = append(parts, fmt.Sprintf("%d", counters[i])) } }
			num := strings.Join(parts, ".")

			indent := strings.Repeat("  ", level-1)
			fmt.Printf("%s%s. %s\n", indent, num, text)
		}
	}

	if len(files) == 0 { process(os.Stdin); return }
	for _, f := range files {
		fh, err := os.Open(f); if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		process(fh); fh.Close()
	}
}
