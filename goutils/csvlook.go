// csvlook - Pretty-print CSV as an ASCII table.
//
// Usage:
//
//	csvlook [OPTIONS] [FILE...]
//	cat data.csv | csvlook
//
// Options:
//
//	-d SEP    Delimiter (default: ,)
//	-n N      Max rows to display (default: all)
//	-c COLS   Only show these columns (comma-separated names or 1-based indices)
//	-s COL    Sort by column name or index
//	-r        Reverse sort
//	-H        No header row
//	-t        Use tab-separated input (shortcut for -d $'\t')
//	-w N      Max column width (truncate, default: 30)
//	-j        JSON output instead of table
//	--style   Border style: ascii|unicode|minimal (default: unicode)
//
// Examples:
//
//	csvlook data.csv
//	cat data.csv | csvlook -n 20 -s age -r
//	csvlook -c name,email users.csv
//	csvlook -t data.tsv
package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	delim    = flag.String("d", ",", "delimiter")
	maxRows  = flag.Int("n", 0, "max rows")
	cols     = flag.String("c", "", "columns")
	sortCol  = flag.String("s", "", "sort column")
	reverse  = flag.Bool("r", false, "reverse sort")
	noHeader = flag.Bool("H", false, "no header")
	tabMode  = flag.Bool("t", false, "TSV mode")
	maxWidth = flag.Int("w", 30, "max column width")
	asJSON   = flag.Bool("j", false, "JSON output")
	style    = flag.String("style", "unicode", "border style")
)

var borders = map[string][9]string{
	"unicode": {"┌", "─", "┬", "├", "─", "┼", "└", "─", "┘"},
	"ascii":   {"+", "-", "+", "+", "-", "+", "+", "-", "+"},
	"minimal": {" ", " ", " ", " ", " ", " ", " ", " ", " "},
}

func truncate(s string, max int) string {
	if max <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}

func colWidth(s string) int {
	return utf8.RuneCountInString(s)
}

func pad(s string, width int) string {
	n := colWidth(s)
	if n >= width {
		return s
	}
	return s + strings.Repeat(" ", width-n)
}

func printTable(headers []string, rows [][]string, bstyle string) {
	b := borders[bstyle]
	if b == ([9]string{}) {
		b = borders["unicode"]
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = colWidth(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				w := colWidth(cell)
				if w > widths[i] {
					widths[i] = w
				}
			}
		}
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	sep := func(left, mid, right, fill string) {
		fmt.Fprint(w, left)
		for i, width := range widths {
			fmt.Fprint(w, strings.Repeat(fill, width+2))
			if i < len(widths)-1 {
				fmt.Fprint(w, mid)
			}
		}
		fmt.Fprintln(w, right)
	}

	printRow := func(cells []string) {
		fmt.Fprint(w, "│")
		for i, width := range widths {
			cell := ""
			if i < len(cells) {
				cell = cells[i]
			}
			fmt.Fprintf(w, " %s │", pad(cell, width))
		}
		fmt.Fprintln(w)
	}

	_ = b

	// Top border
	sep("┌", "┬", "┐", "─")
	// Header
	printRow(headers)
	// Header separator
	sep("├", "┼", "┤", "─")
	// Data rows
	for _, row := range rows {
		printRow(row)
	}
	// Bottom border
	sep("└", "┴", "┘", "─")
	fmt.Fprintf(w, "\n%d rows\n", len(rows))
}

func main() {
	flag.Parse()

	sep := ','
	if *tabMode {
		sep = '\t'
	} else if *delim != "," {
		runes := []rune(*delim)
		if len(runes) > 0 {
			sep = runes[0]
		}
	}

	files := flag.Args()
	var readers []*os.File
	if len(files) == 0 {
		readers = []*os.File{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "csvlook: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}

	var headers []string
	var allRows [][]string

	for _, r := range readers {
		cr := csv.NewReader(bufio.NewReader(r))
		cr.Comma = sep
		cr.LazyQuotes = true
		cr.TrimLeadingSpace = true

		records, err := cr.ReadAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "csvlook: %v\n", err)
			os.Exit(1)
		}
		if len(records) == 0 {
			continue
		}
		if !*noHeader && headers == nil {
			headers = records[0]
			allRows = append(allRows, records[1:]...)
		} else {
			allRows = append(allRows, records...)
		}
	}

	if *noHeader || headers == nil {
		headers = make([]string, 0)
		if len(allRows) > 0 {
			for i := range allRows[0] {
				headers = append(headers, strconv.Itoa(i+1))
			}
		}
	}

	// Filter columns
	colIndices := make([]int, len(headers))
	for i := range headers {
		colIndices[i] = i
	}
	if *cols != "" {
		colIndices = nil
		for _, c := range strings.Split(*cols, ",") {
			c = strings.TrimSpace(c)
			idx, err := strconv.Atoi(c)
			if err == nil {
				colIndices = append(colIndices, idx-1)
			} else {
				for i, h := range headers {
					if h == c {
						colIndices = append(colIndices, i)
						break
					}
				}
			}
		}
	}

	// Select columns
	filterRow := func(row []string) []string {
		out := make([]string, len(colIndices))
		for i, idx := range colIndices {
			if idx < len(row) {
				out[i] = truncate(row[idx], *maxWidth)
			}
		}
		return out
	}

	filteredHeaders := filterRow(headers)
	var filteredRows [][]string
	for _, row := range allRows {
		filteredRows = append(filteredRows, filterRow(row))
	}

	// Sort
	if *sortCol != "" {
		sortIdx := -1
		idx, err := strconv.Atoi(*sortCol)
		if err == nil {
			sortIdx = idx - 1
		} else {
			for i, h := range filteredHeaders {
				if h == *sortCol {
					sortIdx = i
					break
				}
			}
		}
		if sortIdx >= 0 {
			sort.SliceStable(filteredRows, func(i, j int) bool {
				a, b := "", ""
				if sortIdx < len(filteredRows[i]) {
					a = filteredRows[i][sortIdx]
				}
				if sortIdx < len(filteredRows[j]) {
					b = filteredRows[j][sortIdx]
				}
				if *reverse {
					return a > b
				}
				return a < b
			})
		}
	}

	// Limit rows
	if *maxRows > 0 && len(filteredRows) > *maxRows {
		filteredRows = filteredRows[:*maxRows]
	}

	if *asJSON {
		var result []map[string]string
		for _, row := range filteredRows {
			m := make(map[string]string)
			for i, h := range filteredHeaders {
				if i < len(row) {
					m[h] = row[i]
				}
			}
			result = append(result, m)
		}
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))
		return
	}

	printTable(filteredHeaders, filteredRows, *style)
}
