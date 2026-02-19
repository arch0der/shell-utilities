// dirsize - show disk usage of directories, sorted by size
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type entry struct {
	path  string
	size  int64
	files int
}

func humanSize(b int64) string {
	switch {
	case b < 1024: return fmt.Sprintf("%dB", b)
	case b < 1024*1024: return fmt.Sprintf("%.1fK", float64(b)/1024)
	case b < 1024*1024*1024: return fmt.Sprintf("%.1fM", float64(b)/(1024*1024))
	default: return fmt.Sprintf("%.2fG", float64(b)/(1024*1024*1024))
	}
}

func sizeOf(root string) (int64, int) {
	var total int64; count := 0
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil { return nil }
		if !info.IsDir() { total += info.Size(); count++ }
		return nil
	})
	return total, count
}

func main() {
	all := false
	depth := 1
	sort_ := true
	targets := []string{"."}

	args := os.Args[1:]
	rest := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-a", "--all": all = true
		case "--no-sort": sort_ = false
		case "-d": i++; fmt.Sscanf(args[i], "%d", &depth)
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Fprintf(os.Stderr, "dirsize: unknown flag %s\n", args[i])
			} else { rest = append(rest, args[i]) }
		}
	}
	if len(rest) > 0 { targets = rest }

	var entries []entry
	for _, target := range targets {
		des, err := os.ReadDir(target)
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		for _, de := range des {
			if !all && strings.HasPrefix(de.Name(), ".") { continue }
			path := filepath.Join(target, de.Name())
			var sz int64; var fc int
			if de.IsDir() { sz, fc = sizeOf(path) } else {
				info, _ := de.Info()
				if info != nil { sz = info.Size() }; fc = 1
			}
			entries = append(entries, entry{path, sz, fc})
		}
		// also show total
		if len(targets) == 1 {
			total, tc := sizeOf(target)
			_ = depth
			entries = append(entries, entry{target + " (total)", total, tc})
		}
	}

	if sort_ { sort.Slice(entries, func(i, j int) bool { return entries[i].size > entries[j].size }) }
	maxName := 10
	for _, e := range entries { if len(e.path) > maxName { maxName = len(e.path) } }
	if maxName > 60 { maxName = 60 }
	fmt.Printf("%-*s  %8s  %s\n", maxName, "Path", "Size", "Files")
	fmt.Println(strings.Repeat("─", maxName+20))
	for _, e := range entries {
		name := e.path; if len(name) > maxName { name = "..." + name[len(name)-maxName+3:] }
		fmt.Printf("%-*s  %8s  %d\n", maxName, name, humanSize(e.size), e.files)
	}
}
