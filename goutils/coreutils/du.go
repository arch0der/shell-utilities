package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func init() { register("du", runDu) }

func runDu() {
	args := os.Args[1:]
	humanReadable := false
	summarize := false
	all := false
	maxDepth := -1
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--human-readable":
			humanReadable = true
		case a == "-s" || a == "--summarize":
			summarize = true
		case a == "-a" || a == "--all":
			all = true
		case a == "--max-depth" && i+1 < len(args):
			i++
			fmt.Sscan(args[i], &maxDepth)
		case strings.HasPrefix(a, "--max-depth="):
			fmt.Sscan(a[12:], &maxDepth)
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}

	if len(files) == 0 {
		files = []string{"."}
	}

	printSize := func(size int64, path string) {
		if humanReadable {
			fmt.Printf("%s\t%s\n", humanSize(size), path)
		} else {
			fmt.Printf("%d\t%s\n", (size+1023)/1024, path)
		}
	}

	var duDir func(path string, depth int) int64
	duDir = func(path string, depth int) int64 {
		var total int64
		filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if p == path {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			if d.IsDir() {
				subTotal := duDir(p, depth+1)
				if !summarize && (maxDepth < 0 || depth < maxDepth) {
					printSize(subTotal, p)
				}
				total += subTotal
				return filepath.SkipDir
			}
			total += info.Size()
			if all && !summarize {
				printSize(info.Size(), p)
			}
			return nil
		})
		info, _ := os.Lstat(path)
		if info != nil {
			total += info.Size()
		}
		return total
	}

	for _, f := range files {
		info, err := os.Lstat(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "du: %s: %v\n", f, err)
			continue
		}
		var total int64
		if info.IsDir() {
			total = duDir(f, 0)
		} else {
			total = info.Size()
		}
		printSize(total, f)
	}
}
