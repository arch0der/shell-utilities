// treels - Display directory contents as a tree.
//
// Usage:
//
//	treels [OPTIONS] [DIR...]
//
// Options:
//
//	-L N      Max depth (default: unlimited)
//	-a        Show hidden files (starting with .)
//	-d        Directories only
//	-s        Show file sizes
//	-t        Sort by modification time
//	-r        Reverse sort
//	-p        Show permissions
//	-j        JSON output
//	-I PAT    Ignore files matching pattern (glob)
//	--noreport  Omit summary line
//
// Examples:
//
//	treels
//	treels -L 2 /etc
//	treels -a -s ~/projects
//	treels -d -L 3 .
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	maxDepth  = flag.Int("L", -1, "max depth")
	showHide  = flag.Bool("a", false, "show hidden")
	dirsOnly  = flag.Bool("d", false, "dirs only")
	showSize  = flag.Bool("s", false, "show sizes")
	byTime    = flag.Bool("t", false, "sort by time")
	reverse   = flag.Bool("r", false, "reverse sort")
	showPerms = flag.Bool("p", false, "show permissions")
	asJSON    = flag.Bool("j", false, "JSON output")
	ignores   = flag.String("I", "", "ignore pattern")
	noReport  = flag.Bool("noreport", false, "no summary")
)

type Node struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Size     int64   `json:"size,omitempty"`
	Children []*Node `json:"children,omitempty"`
}

var totalFiles, totalDirs int

func humanSize(n int64) string {
	switch {
	case n >= 1<<30:
		return fmt.Sprintf("%.1fG", float64(n)/(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.1fM", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1fK", float64(n)/(1<<10))
	}
	return fmt.Sprintf("%dB", n)
}

func shouldIgnore(name string) bool {
	if !*showHide && strings.HasPrefix(name, ".") {
		return true
	}
	if *ignores != "" {
		matched, _ := filepath.Match(*ignores, name)
		return matched
	}
	return false
}

func printTree(path, prefix string, depth int) {
	if *maxDepth >= 0 && depth > *maxDepth {
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "treels: %v\n", err)
		return
	}

	if *byTime {
		sort.Slice(entries, func(i, j int) bool {
			ii, _ := entries[i].Info()
			ji, _ := entries[j].Info()
			if ii == nil || ji == nil {
				return false
			}
			if *reverse {
				return ii.ModTime().Before(ji.ModTime())
			}
			return ii.ModTime().After(ji.ModTime())
		})
	} else {
		sort.Slice(entries, func(i, j int) bool {
			if *reverse {
				return entries[i].Name() > entries[j].Name()
			}
			return entries[i].Name() < entries[j].Name()
		})
	}

	visible := make([]os.DirEntry, 0, len(entries))
	for _, e := range entries {
		if shouldIgnore(e.Name()) {
			continue
		}
		if *dirsOnly && !e.IsDir() {
			continue
		}
		visible = append(visible, e)
	}

	for i, e := range visible {
		isLast := i == len(visible)-1
		connector := "├── "
		childPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			childPrefix = prefix + "    "
		}

		info, _ := e.Info()
		label := e.Name()

		extras := ""
		if *showPerms && info != nil {
			extras += " [" + info.Mode().String() + "]"
		}
		if *showSize && info != nil && !e.IsDir() {
			extras += " (" + humanSize(info.Size()) + ")"
		}
		if e.IsDir() {
			label += "/"
			totalDirs++
		} else {
			totalFiles++
		}

		fmt.Printf("%s%s%s%s\n", prefix, connector, label, extras)

		if e.IsDir() {
			printTree(filepath.Join(path, e.Name()), childPrefix, depth+1)
		}
	}
}

func buildTree(path string, depth int) *Node {
	info, _ := os.Stat(path)
	node := &Node{Name: filepath.Base(path), Type: "dir"}
	if info != nil && !info.IsDir() {
		node.Type = "file"
		node.Size = info.Size()
		return node
	}

	if *maxDepth >= 0 && depth > *maxDepth {
		return node
	}

	entries, _ := os.ReadDir(path)
	for _, e := range entries {
		if shouldIgnore(e.Name()) {
			continue
		}
		child := buildTree(filepath.Join(path, e.Name()), depth+1)
		node.Children = append(node.Children, child)
	}
	return node
}

func main() {
	flag.Parse()
	dirs := flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	if *asJSON {
		var roots []*Node
		for _, d := range dirs {
			roots = append(roots, buildTree(d, 0))
		}
		var out interface{} = roots
		if len(roots) == 1 {
			out = roots[0]
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return
	}

	for _, d := range dirs {
		fmt.Println(d)
		printTree(d, "", 1)
		fmt.Println()
	}

	if !*noReport {
		fmt.Printf("%d director%s, %d file%s\n",
			totalDirs, plural(totalDirs, "y", "ies"),
			totalFiles, plural(totalFiles, "", "s"))
	}
}

func plural(n int, sing, plur string) string {
	if n == 1 {
		return sing
	}
	return plur
}
