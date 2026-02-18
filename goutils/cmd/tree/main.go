// tree - Display directory structure as a tree
// Usage: tree [dir] [-a] [-L depth]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	all      = flag.Bool("a", false, "Show hidden files")
	maxDepth = flag.Int("L", -1, "Max depth (-1 = unlimited)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: tree [dir] [-a] [-L depth]")
		flag.PrintDefaults()
	}
	flag.Parse()
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}
	fmt.Println(root)
	dirs, files := 0, 0
	printTree(root, "", 0, &dirs, &files)
	fmt.Printf("\n%d directories, %d files\n", dirs, files)
}

func printTree(path, prefix string, depth int, dirs, files *int) {
	if *maxDepth >= 0 && depth >= *maxDepth {
		return
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}
	// Filter hidden if needed
	var visible []os.DirEntry
	for _, e := range entries {
		if !*all && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		visible = append(visible, e)
	}
	for i, entry := range visible {
		connector := "├── "
		childPrefix := prefix + "│   "
		if i == len(visible)-1 {
			connector = "└── "
			childPrefix = prefix + "    "
		}
		fmt.Printf("%s%s%s\n", prefix, connector, entry.Name())
		if entry.IsDir() {
			*dirs++
			printTree(filepath.Join(path, entry.Name()), childPrefix, depth+1, dirs, files)
		} else {
			*files++
		}
	}
}
