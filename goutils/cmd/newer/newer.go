// newer - List files newer than a reference file or timestamp.
//
// Usage:
//
//	newer [OPTIONS] REFERENCE [FILE/DIR...]
//
// Options:
//
//	-r        Recursive directory search
//	-d DATE   Use date string instead of reference file
//	-t        Also print modification times
//	-0        Null-separated output (for xargs)
//	-j        JSON output
//
// Examples:
//
//	newer reference.txt *.go        # files newer than reference.txt
//	newer -d "2024-01-01" /var/log  # files newer than date
//	newer -r -t . reference.txt     # recursive, show times
//	newer -0 ref.txt . | xargs rm
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	recursive = flag.Bool("r", false, "recursive")
	dateStr   = flag.String("d", "", "reference date")
	times     = flag.Bool("t", false, "show mtime")
	null      = flag.Bool("0", false, "null-separated")
	asJSON    = flag.Bool("j", false, "JSON output")
)

type FileEntry struct {
	Path    string    `json:"path"`
	ModTime time.Time `json:"mtime"`
}

var parseFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02",
	"2006/01/02",
	"01/02/2006",
	time.RFC3339,
}

func parseDate(s string) (time.Time, error) {
	for _, f := range parseFormats {
		t, err := time.Parse(f, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse date: %q", s)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 && *dateStr == "" {
		fmt.Fprintln(os.Stderr, "usage: newer REFERENCE [FILE...] or newer -d DATE [FILE...]")
		os.Exit(1)
	}

	var refTime time.Time
	var searchPaths []string

	if *dateStr != "" {
		t, err := parseDate(*dateStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "newer: %v\n", err)
			os.Exit(1)
		}
		refTime = t
		searchPaths = args
	} else {
		refFile := args[0]
		info, err := os.Stat(refFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "newer: %v\n", err)
			os.Exit(1)
		}
		refTime = info.ModTime()
		searchPaths = args[1:]
	}

	if len(searchPaths) == 0 {
		searchPaths = []string{"."}
	}

	var results []FileEntry

	for _, p := range searchPaths {
		info, err := os.Stat(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "newer: %v\n", err)
			continue
		}
		if info.IsDir() {
			walk := func(path string, fi os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if fi.IsDir() && !*recursive && path != p {
					return filepath.SkipDir
				}
				if !fi.IsDir() && fi.ModTime().After(refTime) {
					results = append(results, FileEntry{path, fi.ModTime()})
				}
				return nil
			}
			filepath.Walk(p, walk)
		} else {
			if info.ModTime().After(refTime) {
				results = append(results, FileEntry{p, info.ModTime()})
			}
		}
	}

	if *asJSON {
		b, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(b))
		return
	}

	for _, e := range results {
		if *null {
			fmt.Printf("%s\x00", e.Path)
		} else if *times {
			fmt.Printf("%s\t%s\n", e.ModTime.Format("2006-01-02 15:04:05"), e.Path)
		} else {
			fmt.Println(e.Path)
		}
	}
}
