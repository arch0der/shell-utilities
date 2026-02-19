// tojson - Convert key=value, TSV, or line data to JSON.
//
// Usage:
//
//	tojson [OPTIONS] [FILE...]
//	cat env-file | tojson
//
// Options:
//
//	-f FORMAT  Input format: env|tsv|csv|lines|kv (default: auto)
//	-k KEYS    Comma-separated keys for tsv/lines (e.g. name,age,email)
//	-c         Compact output
//	-a         Output array (one object per line)
//	-t         Try to parse numbers and booleans
//
// Formats:
//   env/kv    KEY=VALUE pairs (one per line)
//   lines     Each line becomes a string in an array
//   tsv       Tab-separated, first line = headers (unless -k)
//   csv       Comma-separated, first line = headers (unless -k)
//
// Examples:
//
//	env | tojson -f env                     # env vars to JSON
//	echo -e "a\tb\tc" | tojson -f tsv -k x,y,z
//	cat /etc/os-release | tojson -f kv
//	ls | tojson -f lines
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	format  = flag.String("f", "auto", "input format")
	keys    = flag.String("k", "", "field keys")
	compact = flag.Bool("c", false, "compact output")
	array   = flag.Bool("a", false, "array mode")
	typed   = flag.Bool("t", false, "parse types")
)

func parseVal(s string) interface{} {
	if !*typed {
		return s
	}
	if s == "true" || s == "yes" {
		return true
	}
	if s == "false" || s == "no" {
		return false
	}
	if s == "null" {
		return nil
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

func marshalOut(v interface{}) {
	var b []byte
	if *compact {
		b, _ = json.Marshal(v)
	} else {
		b, _ = json.MarshalIndent(v, "", "  ")
	}
	fmt.Println(string(b))
}

func processKV(lines []string) {
	result := make(map[string]interface{})
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		idx := strings.Index(l, "=")
		if idx < 0 {
			continue
		}
		k := strings.TrimSpace(l[:idx])
		v := strings.TrimSpace(l[idx+1:])
		v = strings.Trim(v, `"'`)
		result[k] = parseVal(v)
	}
	marshalOut(result)
}

func processLines(lines []string) {
	var result []interface{}
	for _, l := range lines {
		result = append(result, parseVal(l))
	}
	marshalOut(result)
}

func processDelimited(lines []string, delim rune) {
	if len(lines) == 0 {
		marshalOut([]interface{}{})
		return
	}

	fieldKeys := []string{}
	if *keys != "" {
		fieldKeys = strings.Split(*keys, ",")
	}

	startRow := 0
	if len(fieldKeys) == 0 {
		// First line is header
		r := csv.NewReader(strings.NewReader(lines[0]))
		r.Comma = delim
		rec, err := r.Read()
		if err == nil {
			fieldKeys = rec
		}
		startRow = 1
	}

	var result []interface{}
	for _, l := range lines[startRow:] {
		if strings.TrimSpace(l) == "" {
			continue
		}
		r := csv.NewReader(strings.NewReader(l))
		r.Comma = delim
		rec, err := r.Read()
		if err != nil {
			continue
		}
		row := make(map[string]interface{})
		for i, v := range rec {
			k := strconv.Itoa(i)
			if i < len(fieldKeys) {
				k = fieldKeys[i]
			}
			row[k] = parseVal(v)
		}
		if *array {
			marshalOut(row)
		} else {
			result = append(result, row)
		}
	}
	if !*array {
		marshalOut(result)
	}
}

func main() {
	flag.Parse()
	files := flag.Args()

	var allLines []string
	var readers []*os.File
	if len(files) == 0 {
		readers = []*os.File{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "tojson: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}

	for _, r := range readers {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			allLines = append(allLines, sc.Text())
		}
	}

	fmt := *format
	// Auto-detect
	if fmt == "auto" {
		for _, l := range allLines {
			l = strings.TrimSpace(l)
			if l == "" || strings.HasPrefix(l, "#") {
				continue
			}
			if strings.Contains(l, "=") {
				fmt = "kv"
			} else if strings.Contains(l, "\t") {
				fmt = "tsv"
			} else if strings.Contains(l, ",") {
				fmt = "csv"
			} else {
				fmt = "lines"
			}
			break
		}
	}

	switch fmt {
	case "env", "kv":
		processKV(allLines)
	case "tsv":
		processDelimited(allLines, '\t')
	case "csv":
		processDelimited(allLines, ',')
	case "lines":
		processLines(allLines)
	default:
		processKV(allLines)
	}
}
