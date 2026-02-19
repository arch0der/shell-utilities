// table - render delimited data as a formatted ASCII table
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	delim := ","
	header := false
	style := "ascii"
	files := []string{}

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d": i++; delim = args[i]
		case "-t": delim = "\t"
		case "-H", "--header": header = true
		case "--style": i++; style = args[i]
		default: files = append(files, args[i])
		}
	}

	readData := func(r io.Reader) [][]string {
		cr := csv.NewReader(r)
		if delim != "," && len(delim) == 1 { cr.Comma = rune(delim[0]) }
		cr.LazyQuotes = true
		all, err := cr.ReadAll()
		if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		return all
	}

	var rows [][]string
	if len(files) == 0 {
		rows = readData(os.Stdin)
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil { fmt.Fprintln(os.Stderr, err); continue }
			rows = append(rows, readData(fh)...)
			fh.Close()
		}
	}
	if len(rows) == 0 { return }
	_ = bufio.NewScanner // keep import
	_ = style

	// Compute column widths
	cols := 0
	for _, r := range rows { if len(r) > cols { cols = len(r) } }
	widths := make([]int, cols)
	for _, r := range rows {
		for i, cell := range r {
			if len(cell) > widths[i] { widths[i] = len(cell) }
		}
	}

	// Box drawing chars
	var tl, tr, bl, br, h, v, lj, rj, tj, bj, x string
	switch style {
	case "unicode", "box":
		tl="╔"; tr="╗"; bl="╚"; br="╝"; h="═"; v="║"; lj="╠"; rj="╣"; tj="╦"; bj="╩"; x="╬"
	default: // ascii
		tl="+"; tr="+"; bl="+"; br="+"; h="-"; v="|"; lj="+"; rj="+"; tj="+"; bj="+"; x="+"
	}

	sep := func(left, mid, right string) string {
		var b strings.Builder
		b.WriteString(left)
		for i, w := range widths {
			b.WriteString(strings.Repeat(h, w+2))
			if i < cols-1 { b.WriteString(mid) } else { b.WriteString(right) }
		}
		return b.String()
	}

	fmt.Println(sep(tl, tj, tr))
	for ri, row := range rows {
		var b strings.Builder
		b.WriteString(v)
		for i := 0; i < cols; i++ {
			cell := ""; if i < len(row) { cell = row[i] }
			b.WriteString(fmt.Sprintf(" %-*s ", widths[i], cell))
			b.WriteString(v)
		}
		fmt.Println(b.String())
		if ri == 0 && header { fmt.Println(sep(lj, x, rj)) }
	}
	fmt.Println(sep(bl, bj, br))
}
