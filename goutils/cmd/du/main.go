// du - Disk usage
// Usage: du [-h] [-s] [path...]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	human   = flag.Bool("h", false, "Human-readable sizes")
	summary = flag.Bool("s", false, "Display only a total")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: du [-h] [-s] [path...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}
	for _, p := range paths {
		size, err := dirSize(p)
		if err != nil {
			fmt.Fprintln(os.Stderr, "du:", err)
			continue
		}
		if *summary {
			fmt.Printf("%s\t%s\n", formatSize(size), p)
		} else {
			printDu(p)
		}
	}
}

func printDu(path string) {
	info, err := os.Lstat(path)
	if err != nil {
		return
	}
	if info.IsDir() {
		entries, _ := os.ReadDir(path)
		for _, e := range entries {
			printDu(filepath.Join(path, e.Name()))
		}
	}
	size, _ := dirSize(path)
	fmt.Printf("%s\t%s\n", formatSize(size), path)
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func formatSize(b int64) string {
	if !*human {
		return fmt.Sprintf("%d", b/1024)
	}
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1fG", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1fM", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1fK", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%dB", b)
	}
}
