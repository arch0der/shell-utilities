// deltree - safely remove directory trees with confirmation and dry-run
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: deltree [options] <dir> [dir...]
  -f, --force    skip confirmation prompt
  -n, --dry-run  show what would be deleted without deleting
  -v, --verbose  list each file as it's deleted
  -p, --pattern  only delete files matching glob pattern`)
	os.Exit(1)
}

var totalFiles, totalDirs int64
var totalBytes int64

func countTree(root string) (files, dirs int64, bytes int64) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil { return nil }
		if info.IsDir() { dirs++ } else { files++; bytes += info.Size() }
		return nil
	})
	return
}

func humanSize(b int64) string {
	switch {
	case b < 1024: return fmt.Sprintf("%dB", b)
	case b < 1024*1024: return fmt.Sprintf("%.1fKB", float64(b)/1024)
	case b < 1024*1024*1024: return fmt.Sprintf("%.1fMB", float64(b)/(1024*1024))
	default: return fmt.Sprintf("%.2fGB", float64(b)/(1024*1024*1024))
	}
}

func confirm(prompt string) bool {
	fmt.Print(prompt)
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	ans := strings.ToLower(strings.TrimSpace(sc.Text()))
	return ans == "y" || ans == "yes"
}

func main() {
	force, dryRun, verbose := false, false, false
	pattern := ""
	var targets []string

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-f","--force": force = true
		case "-n","--dry-run": dryRun = true
		case "-v","--verbose": verbose = true
		case "-p","--pattern": i++; pattern = os.Args[i]
		default:
			if strings.HasPrefix(os.Args[i], "-") { usage() }
			targets = append(targets, os.Args[i])
		}
	}
	if len(targets) == 0 { usage() }

	for _, root := range targets {
		info, err := os.Stat(root)
		if err != nil { fmt.Fprintf(os.Stderr, "deltree: %v\n", err); continue }
		if !info.IsDir() { fmt.Fprintf(os.Stderr, "deltree: %s is not a directory\n", root); continue }

		files, dirs, bytes := countTree(root)
		fmt.Printf("Target: %s  (%d files, %d dirs, %s)\n", root, files, dirs, humanSize(bytes))

		if dryRun {
			filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil { return nil }
				if pattern != "" {
					matched, _ := filepath.Match(pattern, info.Name())
					if !matched { return nil }
				}
				fmt.Printf("  [dry] %s\n", path)
				return nil
			})
			continue
		}

		if !force {
			if !confirm(fmt.Sprintf("Delete %s? [y/N] ", root)) {
				fmt.Println("Skipped."); continue
			}
		}

		if pattern != "" {
			// Selective delete matching pattern
			filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() { return nil }
				matched, _ := filepath.Match(pattern, info.Name())
				if !matched { return nil }
				if verbose { fmt.Printf("  rm %s\n", path) }
				os.Remove(path)
				return nil
			})
		} else {
			if verbose {
				filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
					if err != nil { return nil }
					fmt.Printf("  rm %s\n", path); return nil
				})
			}
			if err := os.RemoveAll(root); err != nil {
				fmt.Fprintf(os.Stderr, "deltree: %v\n", err); continue
			}
		}
		fmt.Printf("Deleted: %s\n", root)
		totalFiles += files; totalDirs += dirs; totalBytes += bytes
	}
	if len(targets) > 1 {
		fmt.Printf("\nTotal removed: %d files, %d dirs, %s\n", totalFiles, totalDirs, humanSize(totalBytes))
	}
}
