// csvdiff - diff two CSV files by key column, showing added/removed/changed rows
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

const (
	green = "\033[32m"
	red   = "\033[31m"
	yellow= "\033[33m"
	reset = "\033[0m"
)

func readCSV(path string) ([]string, []map[string]string, error) {
	f, err := os.Open(path); if err != nil { return nil, nil, err }
	defer f.Close()
	r := csv.NewReader(f); r.LazyQuotes = true
	recs, err := r.ReadAll(); if err != nil { return nil, nil, err }
	if len(recs) == 0 { return nil, nil, nil }
	headers := recs[0]
	rows := make([]map[string]string, len(recs)-1)
	for i, rec := range recs[1:] {
		m := map[string]string{}
		for j, h := range headers { if j < len(rec) { m[h] = rec[j] } }
		rows[i] = m
	}
	return headers, rows, nil
}

func rowKey(row map[string]string, keys []string) string {
	parts := make([]string, len(keys))
	for i, k := range keys { parts[i] = row[k] }
	return strings.Join(parts, "|")
}

func rowEqual(a, b map[string]string, headers []string) bool {
	for _, h := range headers { if a[h] != b[h] { return false } }
	return true
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: csvdiff [options] <file1.csv> <file2.csv>
  -k <col>    key column(s) for matching rows (comma-separated, default: first col)
  --no-color  disable ANSI color output`)
	os.Exit(1)
}

func main() {
	keySpec := ""
	color := true
	args := os.Args[1:]
	rest := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-k": i++; keySpec = args[i]
		case "--no-color": color = false
		default: rest = append(rest, args[i])
		}
	}
	if len(rest) < 2 { usage() }
	if !color { green[0] = 0; _ = red; _ = yellow; _ = reset }

	h1, rows1, err := readCSV(rest[0]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	h2, rows2, err2 := readCSV(rest[1]); if err2 != nil { fmt.Fprintln(os.Stderr, err2); os.Exit(1) }
	_ = h2

	keys := []string{h1[0]}
	if keySpec != "" { keys = strings.Split(keySpec, ",") }

	idx1 := map[string]map[string]string{}
	for _, r := range rows1 { idx1[rowKey(r, keys)] = r }
	idx2 := map[string]map[string]string{}
	for _, r := range rows2 { idx2[rowKey(r, keys)] = r }

	added, removed, changed := 0, 0, 0
	for k, r := range idx2 {
		if _, ok := idx1[k]; !ok {
			fmt.Printf("%s+ %v%s\n", green, r, reset); added++
		}
	}
	for k, r := range idx1 {
		r2, ok := idx2[k]
		if !ok { fmt.Printf("%s- %v%s\n", red, r, reset); removed++ } else if !rowEqual(r, r2, h1) {
			fmt.Printf("%s~ %v -> %v%s\n", yellow, r, r2, reset); changed++
		}
	}
	fmt.Printf("\nAdded: %d  Removed: %d  Changed: %d\n", added, removed, changed)
}
