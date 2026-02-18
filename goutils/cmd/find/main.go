// find - Search for files
// Usage: find <path> [-name pattern] [-type f|d] [-size +n|-n]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	name    = flag.String("name", "", "File name pattern (e.g. *.go)")
	ftype   = flag.String("type", "", "File type: f=file, d=directory")
	minSize = flag.Int64("size", -1, "Minimum file size in bytes")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: find <path> [-name pattern] [-type f|d] [-size bytes]")
		flag.PrintDefaults()
	}
	flag.Parse()
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintln(os.Stderr, "find:", err)
			return nil
		}
		// Type filter
		if *ftype == "f" && info.IsDir() {
			return nil
		}
		if *ftype == "d" && !info.IsDir() {
			return nil
		}
		// Name filter
		if *name != "" {
			matched, err := filepath.Match(*name, info.Name())
			if err != nil || !matched {
				return nil
			}
		}
		// Size filter
		if *minSize >= 0 && info.Size() < *minSize {
			return nil
		}
		fmt.Println(strings.TrimPrefix(path, "./"))
		return nil
	})
}
