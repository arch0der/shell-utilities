// filesync - sync files between directories (one-way, dry-run safe)
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type syncStats struct{ copied, skipped, deleted int; bytes int64 }

func humanSize(b int64) string {
	if b < 1024 { return fmt.Sprintf("%dB", b) }
	if b < 1024*1024 { return fmt.Sprintf("%.1fK", float64(b)/1024) }
	return fmt.Sprintf("%.1fM", float64(b)/(1024*1024))
}

func copyFile(src, dst string) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil { return 0, err }
	in, err := os.Open(src); if err != nil { return 0, err }
	defer in.Close()
	info, _ := in.Stat()
	out, err := os.Create(dst); if err != nil { return 0, err }
	defer out.Close()
	n, err := io.Copy(out, in); if err != nil { return n, err }
	os.Chtimes(dst, time.Now(), info.ModTime())
	return n, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: filesync [options] <src> <dst>
  -n, --dry-run   show what would be done
  -v, --verbose   list each file
  -d, --delete    delete files in dst not in src
  --exclude <pat> exclude files matching glob pattern (repeatable)`)
	os.Exit(1)
}

func main() {
	dryRun, verbose, delete_ := false, false, false
	var exclude []string
	var args []string

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-n", "--dry-run": dryRun = true
		case "-v", "--verbose": verbose = true
		case "-d", "--delete": delete_ = true
		case "--exclude": i++; exclude = append(exclude, os.Args[i])
		default:
			if strings.HasPrefix(os.Args[i], "-") { usage() }
			args = append(args, os.Args[i])
		}
	}
	if len(args) < 2 { usage() }
	src, dst := args[0], args[1]

	excluded := func(name string) bool {
		for _, p := range exclude {
			if m, _ := filepath.Match(p, name); m { return true }
		}
		return false
	}

	var stats syncStats
	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil { return nil }
		if info.IsDir() || excluded(info.Name()) { return nil }
		rel, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, rel)
		dstInfo, err := os.Stat(dstPath)
		if err == nil && !dstInfo.IsDir() {
			if dstInfo.ModTime().Equal(info.ModTime()) && dstInfo.Size() == info.Size() {
				stats.skipped++; if verbose { fmt.Printf("  = %s\n", rel) }; return nil
			}
		}
		if verbose || dryRun { fmt.Printf("  %s %s\n", func() string { if dryRun { return "[copy]" }; return "→" }(), rel) }
		if !dryRun {
			n, err := copyFile(path, dstPath)
			if err != nil { fmt.Fprintln(os.Stderr, err); return nil }
			stats.bytes += n
		}
		stats.copied++
		return nil
	})

	if delete_ {
		filepath.Walk(dst, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() { return nil }
			rel, _ := filepath.Rel(dst, path)
			srcPath := filepath.Join(src, rel)
			if _, err := os.Stat(srcPath); os.IsNotExist(err) {
				if verbose || dryRun { fmt.Printf("  %s %s\n", func() string { if dryRun { return "[delete]" }; return "✗" }(), rel) }
				if !dryRun { os.Remove(path) }
				stats.deleted++
			}
			return nil
		})
	}

	fmt.Printf("\nSync complete: %d copied (%s), %d skipped, %d deleted\n",
		stats.copied, humanSize(stats.bytes), stats.skipped, stats.deleted)
}
