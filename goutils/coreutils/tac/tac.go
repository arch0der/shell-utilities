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
	files := []string{}
	separator := "\n"
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-s" && i+1 < len(args) {
			i++
			separator = args[i]
		} else if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}
	_ = separator
	tacReader := func(r io.Reader) {
		var lines []string
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
		for i := len(lines) - 1; i >= 0; i-- {
			fmt.Println(lines[i])
		}
	}
	if len(files) == 0 {
		tacReader(os.Stdin)
		return
	}
	for _, f := range files {
		fh, _ := os.Open(f)
		tacReader(fh)
		fh.Close()
	}
}
