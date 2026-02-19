// treeprint - render a directory tree (or indented text) as Unicode tree
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	branch = "├── "
	last   = "└── "
	pipe   = "│   "
	space  = "    "
)

func walk(dir, prefix string, depth, maxDepth int, showHidden bool, stats *[2]int) {
	if maxDepth > 0 && depth > maxDepth { return }
	entries, err := os.ReadDir(dir)
	if err != nil { fmt.Fprintf(os.Stderr, "treeprint: %v\n", err); return }
	if !showHidden {
		var visible []os.DirEntry
		for _, e := range entries { if !strings.HasPrefix(e.Name(), ".") { visible = append(visible, e) } }
		entries = visible
	}
	sort.Slice(entries, func(i, j int) bool {
		ai, aj := entries[i].IsDir(), entries[j].IsDir()
		if ai != aj { return ai }
		return entries[i].Name() < entries[j].Name()
	})
	for i, e := range entries {
		isLast := i == len(entries)-1
		connector := branch; if isLast { connector = last }
		if e.IsDir() {
			fmt.Printf("%s%s%s/\n", prefix, connector, e.Name())
			stats[0]++
			newPrefix := prefix + pipe; if isLast { newPrefix = prefix + space }
			walk(filepath.Join(dir, e.Name()), newPrefix, depth+1, maxDepth, showHidden, stats)
		} else {
			fmt.Printf("%s%s%s\n", prefix, connector, e.Name())
			stats[1]++
		}
	}
}

func fromStdin() {
	sc := bufio.NewScanner(os.Stdin)
	var lines []string
	for sc.Scan() { lines = append(lines, sc.Text()) }
	// Print as tree based on indentation
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		depth := (len(line) - len(trimmed)) / 2
		prefix := strings.Repeat(space, depth)
		fmt.Printf("%s%s%s\n", prefix, branch, trimmed)
	}
}

func main() {
	maxDepth := 0
	showHidden := false
	dir := "."
	stdin := false

	for _, arg := range os.Args[1:] {
		switch {
		case arg == "-a": showHidden = true
		case arg == "--stdin": stdin = true
		case strings.HasPrefix(arg, "-L"):
			fmt.Sscanf(arg[2:], "%d", &maxDepth)
		default:
			if strings.HasPrefix(arg, "-") { break }
			dir = arg
		}
	}

	if stdin { fromStdin(); return }

	info, err := os.Stat(dir)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	_ = info

	fmt.Printf("%s\n", dir)
	var stats [2]int
	walk(dir, "", 1, maxDepth, showHidden, &stats)
	fmt.Printf("\n%d directories, %d files\n", stats[0], stats[1])
}
