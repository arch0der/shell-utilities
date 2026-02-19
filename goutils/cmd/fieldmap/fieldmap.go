// fieldmap - rearrange, rename, or filter delimited fields
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: fieldmap [options] <field_spec,...>
  -d <delim>   input delimiter (default: auto-detect CSV/TSV)
  -o <delim>   output delimiter (default: same as input)
  -H           first line is header, use names in spec
  field_spec:  1,3,2  (reorder by index, 1-based)  OR  name1,name2  (with -H)
  Use 0 to emit empty field.`)
	os.Exit(1)
}

func main() {
	inDelim := ","
	outDelim := ""
	hasHeader := false
	var spec string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d": i++; inDelim = args[i]
		case "-o": i++; outDelim = args[i]
		case "-H": hasHeader = true
		default:
			if spec == "" { spec = args[i] } else { usage() }
		}
	}
	if spec == "" { usage() }
	if outDelim == "" { outDelim = inDelim }

	r := csv.NewReader(os.Stdin)
	if inDelim != "," && len(inDelim) == 1 { r.Comma = rune(inDelim[0]) }
	r.LazyQuotes = true; r.TrimLeadingSpace = true

	var header []string
	if hasHeader {
		var err error
		header, err = r.Read()
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	}

	// Parse spec
	rawSpec := strings.Split(spec, ",")
	indices := make([]int, len(rawSpec)) // -1 = empty
	for i, s := range rawSpec {
		s = strings.TrimSpace(s)
		if s == "0" || s == "" { indices[i] = -1; continue }
		if n, err := strconv.Atoi(s); err == nil { indices[i] = n - 1; continue }
		// Name lookup
		found := false
		for j, h := range header {
			if h == s { indices[i] = j; found = true; break }
		}
		if !found { fmt.Fprintf(os.Stderr, "fieldmap: unknown field %q\n", s); os.Exit(1) }
	}

	w := csv.NewWriter(os.Stdout)
	if outDelim != "," && len(outDelim) == 1 { w.Comma = rune(outDelim[0]) }

	write := func(row []string) {
		out := make([]string, len(indices))
		for i, idx := range indices {
			if idx < 0 || idx >= len(row) { out[i] = "" } else { out[i] = row[idx] }
		}
		w.Write(out)
	}
	if hasHeader { write(header) }

	for {
		row, err := r.Read()
		if err == io.EOF { break }
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		write(row)
	}
	w.Flush()
}
