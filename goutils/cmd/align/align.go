// align - align columns in delimited text for pretty printing
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	delim := "\t"
	padding := 2
	right := false
	sep := "  "

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d": i++; delim = args[i]
		case "-p": i++; fmt.Sscanf(args[i], "%d", &padding)
		case "-r": right = true
		case "-s": i++; sep = args[i]
		}
	}

	sc := bufio.NewScanner(os.Stdin)
	var rows [][]string
	for sc.Scan() {
		line := sc.Text()
		if delim == "  " || delim == "ws" {
			rows = append(rows, strings.Fields(line))
		} else {
			rows = append(rows, strings.Split(line, delim))
		}
	}

	if len(rows) == 0 { return }
	cols := 0
	for _, r := range rows { if len(r) > cols { cols = len(r) } }
	widths := make([]int, cols)
	for _, r := range rows {
		for i, cell := range r {
			if len(cell) > widths[i] { widths[i] = len(cell) }
		}
	}
	pad := strings.Repeat(" ", padding)
	for _, r := range rows {
		var parts []string
		for i := 0; i < cols; i++ {
			cell := ""; if i < len(r) { cell = r[i] }
			w := widths[i] + padding
			if right {
				parts = append(parts, fmt.Sprintf("%*s", w, cell))
			} else {
				parts = append(parts, fmt.Sprintf("%-*s", w, cell))
			}
		}
		fmt.Println(strings.TrimRight(strings.Join(parts, sep), " ") + pad[:0])
	}
	_ = pad
}
