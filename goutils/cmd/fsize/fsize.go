// fsize - Display file sizes in human-readable format.
//
// Usage:
//
//	fsize [OPTIONS] [FILE...]
//
// Options:
//
//	-b        Bytes (raw)
//	-k        Kilobytes
//	-m        Megabytes
//	-g        Gigabytes
//	-si       SI units (1K=1000, default: IEC 1K=1024)
//	-s        Sort by size (descending)
//	-t        Show total
//	-j        JSON output
//
// Examples:
//
//	fsize *.log
//	fsize -s ~/Downloads/*
//	fsize -t /var/log/*
//	du -b * | fsize -    # read sizes from stdin
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	rawBytes = flag.Bool("b", false, "raw bytes")
	kb       = flag.Bool("k", false, "kilobytes")
	mb       = flag.Bool("m", false, "megabytes")
	gb       = flag.Bool("g", false, "gigabytes")
	si       = flag.Bool("si", false, "SI units (base 1000)")
	sortSize = flag.Bool("s", false, "sort by size")
	total    = flag.Bool("t", false, "show total")
	asJSON   = flag.Bool("j", false, "JSON output")
)

type Entry struct {
	Name string `json:"name"`
	Size int64  `json:"bytes"`
	Fmt  string `json:"formatted"`
}

func fmtSize(n int64) string {
	if *rawBytes {
		return fmt.Sprintf("%d", n)
	}
	base := int64(1024)
	if *si {
		base = 1000
	}
	suffixes := []string{"B", "K", "M", "G", "T", "P"}
	if *si {
		suffixes = []string{"B", "KB", "MB", "GB", "TB", "PB"}
	}
	if *kb {
		return fmt.Sprintf("%.1fK", float64(n)/float64(base))
	}
	if *mb {
		return fmt.Sprintf("%.1fM", float64(n)/float64(base*base))
	}
	if *gb {
		return fmt.Sprintf("%.1fG", float64(n)/float64(base*base*base))
	}
	// auto
	f := float64(n)
	i := 0
	for f >= float64(base) && i < len(suffixes)-1 {
		f /= float64(base)
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%dB", n)
	}
	return fmt.Sprintf("%.1f%s", f, suffixes[i])
}

func main() {
	flag.Parse()
	files := flag.Args()

	var entries []Entry

	// Support reading from stdin (e.g. du -b * | fsize -)
	if len(files) == 1 && files[0] == "-" {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			parts := strings.Fields(sc.Text())
			if len(parts) < 2 {
				continue
			}
			sz, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				continue
			}
			name := strings.Join(parts[1:], " ")
			entries = append(entries, Entry{Name: name, Size: sz, Fmt: fmtSize(sz)})
		}
	} else if len(files) == 0 {
		// current directory
		dir, _ := os.ReadDir(".")
		for _, d := range dir {
			info, err := d.Info()
			if err != nil {
				continue
			}
			entries = append(entries, Entry{Name: d.Name(), Size: info.Size(), Fmt: fmtSize(info.Size())})
		}
	} else {
		for _, f := range files {
			info, err := os.Stat(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fsize: %v\n", err)
				continue
			}
			entries = append(entries, Entry{Name: f, Size: info.Size(), Fmt: fmtSize(info.Size())})
		}
	}

	if *sortSize {
		sort.Slice(entries, func(i, j int) bool { return entries[i].Size > entries[j].Size })
	}

	if *asJSON {
		b, _ := json.MarshalIndent(entries, "", "  ")
		fmt.Println(string(b))
		return
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	var totalBytes int64
	for _, e := range entries {
		fmt.Fprintf(w, "%8s  %s\n", e.Fmt, e.Name)
		totalBytes += e.Size
	}
	if *total {
		fmt.Fprintf(w, "%8s  total\n", fmtSize(totalBytes))
	}
}
