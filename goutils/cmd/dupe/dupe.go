// dupe - Find duplicate files by content hash.
//
// Usage:
//
//	dupe [OPTIONS] DIR [DIR...]
//
// Options:
//
//	-r         Recursive (default: true)
//	--no-r     Non-recursive
//	-d         Delete duplicates (keeps first found, interactive)
//	-f         Force delete without confirmation (use with -d)
//	-min SIZE  Minimum file size to consider (default: 1)
//	-j         JSON output
//	-0         Print filenames separated by null (for xargs -0)
//
// Examples:
//
//	dupe ~/Downloads
//	dupe -min 1M ~/Photos ~/Backup
//	dupe -d ~/Downloads              # interactively delete dupes
//	dupe -j . | jq '.[] | .paths'
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	recursive = flag.Bool("r", true, "recursive")
	doDelete  = flag.Bool("d", false, "delete duplicates")
	force     = flag.Bool("f", false, "force delete")
	minSize   = flag.String("min", "1", "minimum file size")
	asJSON    = flag.Bool("j", false, "JSON output")
	null      = flag.Bool("0", false, "null-separated output")
)

func parseSize(s string) int64 {
	s = strings.ToUpper(strings.TrimSpace(s))
	mults := map[string]int64{"K": 1024, "M": 1024 * 1024, "G": 1024 * 1024 * 1024}
	for suffix, mult := range mults {
		if strings.HasSuffix(s, suffix) {
			var n int64
			fmt.Sscanf(s[:len(s)-1], "%d", &n)
			return n * mult
		}
	}
	var n int64
	fmt.Sscanf(s, "%d", &n)
	return n
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

type DupeGroup struct {
	Hash  string   `json:"hash"`
	Size  int64    `json:"size"`
	Paths []string `json:"paths"`
}

func main() {
	flag.Parse()
	dirs := flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	minBytes := parseSize(*minSize)
	hashes := make(map[string][]string)
	sizes := make(map[string]int64)

	for _, dir := range dirs {
		walk := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "dupe: %v\n", err)
				return nil
			}
			if info.IsDir() {
				if !*recursive && path != dir {
					return filepath.SkipDir
				}
				return nil
			}
			if !info.Mode().IsRegular() || info.Size() < minBytes {
				return nil
			}
			hash, err := hashFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dupe: %v\n", err)
				return nil
			}
			hashes[hash] = append(hashes[hash], path)
			sizes[hash] = info.Size()
			return nil
		}
		filepath.Walk(dir, walk)
	}

	var groups []DupeGroup
	for hash, paths := range hashes {
		if len(paths) < 2 {
			continue
		}
		sort.Strings(paths)
		groups = append(groups, DupeGroup{Hash: hash, Size: sizes[hash], Paths: paths})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Size > groups[j].Size })

	if len(groups) == 0 {
		if !*asJSON {
			fmt.Println("No duplicates found.")
		} else {
			fmt.Println("[]")
		}
		return
	}

	if *asJSON {
		b, _ := json.MarshalIndent(groups, "", "  ")
		fmt.Println(string(b))
		return
	}

	for _, g := range groups {
		if *null {
			for _, p := range g.Paths[1:] {
				fmt.Printf("%s\x00", p)
			}
			continue
		}
		fmt.Printf("--- %s (%d bytes) ---\n", g.Hash[:12], g.Size)
		for i, p := range g.Paths {
			marker := " "
			if i == 0 {
				marker = "*"
			}
			fmt.Printf("  %s %s\n", marker, p)
		}

		if *doDelete {
			for _, p := range g.Paths[1:] {
				if !*force {
					fmt.Printf("Delete %s? [y/N] ", p)
					var ans string
					fmt.Scanln(&ans)
					if strings.ToLower(ans) != "y" {
						continue
					}
				}
				if err := os.Remove(p); err != nil {
					fmt.Fprintf(os.Stderr, "dupe: %v\n", err)
				} else {
					fmt.Printf("  deleted: %s\n", p)
				}
			}
		}
		fmt.Println()
	}
}
