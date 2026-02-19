package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	pageLength := 66
	pageWidth := 72
	columns := 1
	header := ""
	noHeader := false
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-l" && i+1 < len(args):
			i++
			fmt.Sscan(args[i], &pageLength)
		case strings.HasPrefix(a, "-l"):
			fmt.Sscan(a[2:], &pageLength)
		case a == "-w" && i+1 < len(args):
			i++
			fmt.Sscan(args[i], &pageWidth)
		case a == "-h" && i+1 < len(args):
			i++
			header = args[i]
		case a == "-t":
			noHeader = true
		case len(a) > 1 && a[0] == '-':
			fmt.Sscan(a[1:], &columns)
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}
	_ = pageWidth

	prFile := func(r io.Reader, name string) {
		var lines []string
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
		now := time.Now().Format("2006-01-02 15:04")
		pageNum := 1
		linesPerPage := pageLength - 5 // header takes 5 lines
		if noHeader {
			linesPerPage = pageLength
		}

		for i := 0; i < len(lines); i += linesPerPage * columns {
			if !noHeader {
				title := header
				if title == "" {
					title = name
				}
				fmt.Printf("\n\n%s  %s  Page %d\n\n\n", now, title, pageNum)
			}
			for row := 0; row < linesPerPage; row++ {
				lineIdx := i + row
				if lineIdx >= len(lines) {
					break
				}
				if columns > 1 {
					var cols []string
					for c := 0; c < columns; c++ {
						idx := i + c*linesPerPage + row
						if idx < len(lines) {
							cols = append(cols, lines[idx])
						}
					}
					fmt.Println(strings.Join(cols, "\t"))
				} else {
					fmt.Println(lines[lineIdx])
				}
			}
			pageNum++
		}
	}

	if len(files) == 0 {
		prFile(os.Stdin, "stdin")
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pr: %s: %v\n", f, err)
			continue
		}
		prFile(fh, f)
		fh.Close()
	}
}
