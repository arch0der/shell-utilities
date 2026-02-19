package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	tabWidth := 8
	files := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-t" && i+1 < len(args) {
			i++
			fmt.Sscan(args[i], &tabWidth)
		} else if strings.HasPrefix(a, "-t") {
			fmt.Sscan(a[2:], &tabWidth)
		} else if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}
	process := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			line := sc.Text()
			var out strings.Builder
			col := 0
			for _, c := range line {
				if c == '\t' {
					spaces := tabWidth - (col % tabWidth)
					for i := 0; i < spaces; i++ {
						out.WriteByte(' ')
					}
					col += spaces
				} else {
					out.WriteRune(c)
					col++
				}
			}
			fmt.Println(out.String())
		}
	}
	if len(files) == 0 {
		process(os.Stdin)
		return
	}
	for _, f := range files {
		fh, _ := os.Open(f)
		process(fh)
		fh.Close()
	}
}
