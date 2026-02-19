// cols - Select, reorder, and filter columns from delimited text.
//
// Usage:
//
//	cols [OPTIONS] COLS [FILE...]
//	cat data.csv | cols 1,3,2
//
// Options:
//
//	-d SEP    Input delimiter (default: whitespace)
//	-o SEP    Output delimiter (default: same as input, or tab for whitespace)
//	-H        First line is header; select by name
//	-n        Print column names/indices from header line
//	-r        Reverse selection (print all except specified cols)
//	-t        Trim whitespace from each field
//
// Column spec: comma-separated 1-based indices or names (with -H).
//   Range: 2-5 means cols 2,3,4,5
//   Negative: -1 means last column
//
// Examples:
//
//	cat /etc/passwd | cols -d: 1,3,7    # name, uid, shell
//	ps aux | cols 1,2,11               # user, pid, command
//	cat data.csv | cols -H name,email   # by column name
//	cat file.tsv | cols -d $'\t' 2-4
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	inDelim  = flag.String("d", "", "input delimiter")
	outDelim = flag.String("o", "", "output delimiter")
	header   = flag.Bool("H", false, "first line is header")
	names    = flag.Bool("n", false, "print column names")
	reverse  = flag.Bool("r", false, "reverse selection")
	trim     = flag.Bool("t", false, "trim fields")
)

func splitLine(line, delim string) []string {
	if delim == "" {
		return strings.Fields(line)
	}
	return strings.Split(line, delim)
}

func parseColSpec(spec string, headers []string) ([]int, error) {
	var indices []int
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// Range: 2-5
		if strings.Contains(part, "-") && !strings.HasPrefix(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			start, err1 := strconv.Atoi(rangeParts[0])
			end, err2 := strconv.Atoi(rangeParts[1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			for i := start; i <= end; i++ {
				indices = append(indices, i-1) // 0-based
			}
			continue
		}
		// Named column
		if n, err := strconv.Atoi(part); err != nil {
			if len(headers) == 0 {
				return nil, fmt.Errorf("column name %q requires -H flag", part)
			}
			found := false
			for i, h := range headers {
				if h == part {
					indices = append(indices, i)
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("column %q not found", part)
			}
		} else {
			if n < 0 {
				indices = append(indices, n) // negative handled at print time
			} else {
				indices = append(indices, n-1)
			}
		}
	}
	return indices, nil
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: cols [OPTIONS] COLS [FILE...]")
		os.Exit(1)
	}

	colSpec := args[0]
	files := args[1:]

	outSep := *outDelim
	if outSep == "" {
		if *inDelim == "" {
			outSep = "\t"
		} else {
			outSep = *inDelim
		}
	}

	var readers []*os.File
	if len(files) == 0 {
		readers = []*os.File{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cols: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	var headerRow []string
	var indices []int
	first := true

	for _, r := range readers {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1024*1024), 1024*1024)
		for sc.Scan() {
			line := sc.Text()
			fields := splitLine(line, *inDelim)
			if *trim {
				for i, f := range fields {
					fields[i] = strings.TrimSpace(f)
				}
			}

			if first && *header {
				headerRow = fields
				first = false
				if *names {
					for i, h := range headerRow {
						fmt.Fprintf(w, "%d: %s\n", i+1, h)
					}
					os.Exit(0)
				}
				var err error
				indices, err = parseColSpec(colSpec, headerRow)
				if err != nil {
					fmt.Fprintf(os.Stderr, "cols: %v\n", err)
					os.Exit(1)
				}
				// still print header
			} else if first {
				first = false
				var err error
				indices, err = parseColSpec(colSpec, nil)
				if err != nil {
					fmt.Fprintf(os.Stderr, "cols: %v\n", err)
					os.Exit(1)
				}
			}

			n := len(fields)
			var selected []string
			if *reverse {
				skip := make(map[int]bool)
				for _, idx := range indices {
					if idx < 0 {
						idx = n + idx
					}
					skip[idx] = true
				}
				for i, f := range fields {
					if !skip[i] {
						selected = append(selected, f)
					}
				}
			} else {
				for _, idx := range indices {
					if idx < 0 {
						idx = n + idx
					}
					if idx >= 0 && idx < n {
						selected = append(selected, fields[idx])
					} else {
						selected = append(selected, "")
					}
				}
			}
			fmt.Fprintln(w, strings.Join(selected, outSep))
		}
	}
}
