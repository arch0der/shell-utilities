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
	width := 80
	breakSpaces := false
	files := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-w" && i+1 < len(args):
			i++
			fmt.Sscan(args[i], &width)
		case strings.HasPrefix(a, "-w"):
			fmt.Sscan(a[2:], &width)
		case a == "-s" || a == "--spaces":
			breakSpaces = true
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}
	foldLine := func(line string) {
		for len(line) > width {
			pos := width
			if breakSpaces {
				for pos > 0 && line[pos-1] != ' ' {
					pos--
				}
				if pos == 0 {
					pos = width
				}
			}
			fmt.Println(line[:pos])
			line = line[pos:]
		}
		fmt.Println(line)
	}
	foldReader := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for sc.Scan() {
			foldLine(sc.Text())
		}
	}
	if len(files) == 0 {
		foldReader(os.Stdin)
		return
	}
	for _, f := range files {
		fh, _ := os.Open(f)
		foldReader(fh)
		fh.Close()
	}
}
