package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func init() { register("fmt", runFmt) }

func runFmt() {
	args := os.Args[1:]
	width := 75
	files := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-w" && i+1 < len(args) {
			i++
			fmt.Sscan(args[i], &width)
		} else if strings.HasPrefix(a, "-w") {
			fmt.Sscan(a[2:], &width)
		} else if len(a) > 1 && a[0] == '-' {
			fmt.Sscan(a[1:], &width)
		} else {
			files = append(files, a)
		}
	}
	fmtReader := func(r io.Reader) {
		scanner := bufio.NewScanner(r)
		var para []string
		flush := func() {
			if len(para) == 0 {
				return
			}
			words := strings.Fields(strings.Join(para, " "))
			line := ""
			for _, w := range words {
				if line == "" {
					line = w
				} else if len(line)+1+len(w) <= width {
					line += " " + w
				} else {
					fmt.Println(line)
					line = w
				}
			}
			if line != "" {
				fmt.Println(line)
			}
			para = para[:0]
		}
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				flush()
				fmt.Println()
			} else {
				para = append(para, text)
			}
		}
		flush()
	}
	if len(files) == 0 {
		fmtReader(os.Stdin)
		return
	}
	for _, f := range files {
		fh, _ := os.Open(f)
		fmtReader(fh)
		fh.Close()
	}
}
