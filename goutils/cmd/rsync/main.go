// rsync - Synchronize files between directories (local mode)
// Usage: rsync [-r] [-a] [-v] [-n] [-u] [--delete] src dst
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	recursive = flag.Bool("r", false, "Recursive")
	archive   = flag.Bool("a", false, "Archive mode (-r + preserve attrs)")
	verbose   = flag.Bool("v", false, "Verbose output")
	dryRun    = flag.Bool("n", false, "Dry run (no changes)")
	update    = flag.Bool("u", false, "Skip if destination is newer")
	delete    = flag.Bool("delete", false, "Delete extra files from destination")
	progress  = flag.Bool("progress", false, "Show progress")
	exclude   = flag.String("exclude", "", "Exclude pattern")
)

var stats struct {
	files, dirs, bytes int64
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: rsync [-ravnu] [--delete] [--exclude=PATTERN] src dst")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if *archive {
		*recursive = true
	}

	src := flag.Arg(0)
	dst := flag.Arg(1)

	// If src ends with /, copy contents; otherwise copy directory itself
	srcInfo, err := os.Stat(src)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rsync:", err)
		os.Exit(1)
	}

	if srcInfo.IsDir() {
		if *recursive || *archive {
			if err := syncDir(src, dst); err != nil {
				fmt.Fprintln(os.Stderr, "rsync:", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "rsync: omitting directory (use -r or -a)")
			os.Exit(1)
		}
	} else {
		if err := syncFile(src, dst); err != nil {
			fmt.Fprintln(os.Stderr, "rsync:", err)
			os.Exit(1)
		}
	}

	if *verbose || *dryRun {
		fmt.Printf("\nrsync done: %d files, %d dirs, %s transferred\n",
			stats.files, stats.dirs, humanBytes(stats.bytes))
	}
}

func shouldExclude(path string) bool {
	if *exclude == "" {
		return false
	}
	matched, _ := filepath.Match(*exclude, filepath.Base(path))
	return matched || strings.Contains(path, *exclude)
}

func syncDir(src, dst string) error {
	if shouldExclude(src) {
		return nil
	}

	if !*dryRun {
		srcInfo, err := os.Stat(src)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
			return err
		}
	}
	if *verbose {
		fmt.Printf("created directory %s\n", dst)
	}
	stats.dirs++

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Track destination files for --delete
	dstFiles := map[string]bool{}
	if *delete {
		if dstEntries, err := os.ReadDir(dst); err == nil {
			for _, e := range dstEntries {
				dstFiles[e.Name()] = true
			}
		}
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		delete(dstFiles, entry.Name())

		if shouldExclude(srcPath) {
			continue
		}

		if entry.IsDir() {
			if *recursive || *archive {
				if err := syncDir(srcPath, dstPath); err != nil {
					return err
				}
			}
		} else {
			if err := syncFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	// Delete extra files in destination
	if *delete {
		for name := range dstFiles {
			dstPath := filepath.Join(dst, name)
			if *verbose {
				fmt.Printf("deleting %s\n", dstPath)
			}
			if !*dryRun {
				os.RemoveAll(dstPath)
			}
		}
	}
	return nil
}

func syncFile(src, dst string) error {
	if shouldExclude(src) {
		return nil
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Check if update needed
	if *update {
		dstInfo, err := os.Stat(dst)
		if err == nil && dstInfo.ModTime().After(srcInfo.ModTime()) {
			return nil
		}
	}

	// Check if files are identical (same size + mtime)
	if dstInfo, err := os.Stat(dst); err == nil {
		if dstInfo.Size() == srcInfo.Size() &&
			dstInfo.ModTime().Equal(srcInfo.ModTime()) {
			return nil
		}
	}

	if *verbose || *progress {
		fmt.Printf("%s\n", src)
	}
	if *dryRun {
		stats.files++
		stats.bytes += srcInfo.Size()
		return nil
	}

	// Handle dst as directory
	if dstInfo, err := os.Stat(dst); err == nil && dstInfo.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	if err := copyFile(src, dst, srcInfo); err != nil {
		return err
	}
	stats.files++
	stats.bytes += srcInfo.Size()
	return nil
}

func copyFile(src, dst string, srcInfo os.FileInfo) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		// Try creating parent dir
		os.MkdirAll(filepath.Dir(dst), 0755)
		df, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
		if err != nil {
			return err
		}
	}

	if _, err := io.Copy(df, sf); err != nil {
		df.Close()
		return err
	}
	df.Close()

	// Preserve modification time
	if *archive {
		os.Chtimes(dst, time.Now(), srcInfo.ModTime())
	}
	return nil
}

func humanBytes(n int64) string {
	switch {
	case n >= 1<<30:
		return fmt.Sprintf("%.1fGB", float64(n)/(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.1fMB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1fKB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%dB", n)
	}
}
