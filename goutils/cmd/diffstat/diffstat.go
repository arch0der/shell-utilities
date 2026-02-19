// diffstat - parse unified diff output and show a summary bar chart
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type fileStat struct {
	name    string
	added   int
	removed int
}

func parseDiff(r io.Reader) []fileStat {
	sc := bufio.NewScanner(r)
	var stats []fileStat
	var cur *fileStat
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ "):
			if strings.HasPrefix(line, "+++ ") {
				name := strings.TrimPrefix(line, "+++ ")
				name = strings.TrimPrefix(name, "b/")
				if strings.Contains(name, "\t") { name = name[:strings.Index(name, "\t")] }
				if name != "/dev/null" {
					stats = append(stats, fileStat{name: name})
					cur = &stats[len(stats)-1]
				} else { cur = nil }
			}
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			if cur != nil { cur.added++ }
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			if cur != nil { cur.removed++ }
		}
	}
	return stats
}

func main() {
	barWidth := 50
	var r io.Reader = os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		defer f.Close(); r = f
	}
	stats := parseDiff(r)
	if len(stats) == 0 { fmt.Println("No diff data found."); return }

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].added+stats[i].removed > stats[j].added+stats[j].removed
	})

	maxChange := 1
	for _, s := range stats {
		if s.added+s.removed > maxChange { maxChange = s.added + s.removed }
	}

	totalAdded, totalRemoved := 0, 0
	nameWidth := 20
	for _, s := range stats {
		if len(s.name) > nameWidth { nameWidth = len(s.name) }
		totalAdded += s.added; totalRemoved += s.removed
	}
	if nameWidth > 50 { nameWidth = 50 }

	fmt.Printf(" %-*s │ %-6s %-6s %s\n", nameWidth, "File", "+", "-", "Graph")
	fmt.Println(strings.Repeat("─", nameWidth+2) + "┼" + strings.Repeat("─", barWidth+20))
	for _, s := range stats {
		total := s.added + s.removed
		addBars := 0; remBars := 0
		if total > 0 {
			addBars = s.added * barWidth / maxChange
			remBars = s.removed * barWidth / maxChange
			if addBars == 0 && s.added > 0 { addBars = 1 }
			if remBars == 0 && s.removed > 0 { remBars = 1 }
		}
		name := s.name
		if len(name) > nameWidth { name = "..." + name[len(name)-nameWidth+3:] }
		bar := strings.Repeat("+", addBars) + strings.Repeat("-", remBars)
		fmt.Printf(" %-*s │ %-6d %-6d %s\n", nameWidth, name, s.added, s.removed, bar)
	}
	fmt.Println(strings.Repeat("─", nameWidth+2) + "┴" + strings.Repeat("─", barWidth+20))
	fmt.Printf(" %d files changed, %d insertions(+), %d deletions(-)\n", len(stats), totalAdded, totalRemoved)
}
