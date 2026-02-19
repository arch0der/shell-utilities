package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	bodyStyle := "t" // t=non-empty, a=all, n=none
	width := 6
	sep := "\t"
	startNum := 1
	increment := 1
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-b" && i+1 < len(args):
			i++
			bodyStyle = args[i]
		case strings.HasPrefix(a, "-b"):
			bodyStyle = a[2:]
		case a == "-w" && i+1 < len(args):
			i++
			width, _ = strconv.Atoi(args[i])
		case strings.HasPrefix(a, "-w"):
			width, _ = strconv.Atoi(a[2:])
		case a == "-s" && i+1 < len(args):
			i++
			sep = args[i]
		case strings.HasPrefix(a, "-s"):
			sep = a[2:]
		case a == "-v" && i+1 < len(args):
			i++
			startNum, _ = strconv.Atoi(args[i])
		case a == "-i" && i+1 < len(args):
			i++
			increment, _ = strconv.Atoi(args[i])
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}

	nlReader := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		n := startNum
		for sc.Scan() {
			line := sc.Text()
			shouldNum := false
			switch bodyStyle {
			case "a":
				shouldNum = true
			case "t":
				shouldNum = strings.TrimSpace(line) != ""
			case "n":
				shouldNum = false
			}
			if shouldNum {
				fmt.Printf("%*d%s%s\n", width, n, sep, line)
				n += increment
			} else {
				fmt.Printf("%s%s\n", strings.Repeat(" ", width+len(sep)), line)
			}
		}
	}

	if len(files) == 0 {
		nlReader(os.Stdin)
		return
	}
	for _, f := range files {
		fh, _ := os.Open(f)
		nlReader(fh)
		fh.Close()
	}
}
