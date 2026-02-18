// tail - Print last N lines of a file
// Usage: tail [-n lines] [-f] [file]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	n      = flag.Int("n", 10, "Number of lines to print")
	follow = flag.Bool("f", false, "Follow file as it grows")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: tail [-n lines] [-f] [file]")
		flag.PrintDefaults()
	}
	flag.Parse()

	var f *os.File
	var err error
	if flag.NArg() == 0 {
		f = os.Stdin
	} else {
		f, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "tail:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	lines := readLastN(f, *n)
	for _, l := range lines {
		fmt.Println(l)
	}

	if *follow {
		reader := bufio.NewReader(f)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			if err != nil {
				break
			}
			fmt.Print(line)
		}
	}
}

func readLastN(f *os.File, n int) []string {
	scanner := bufio.NewScanner(f)
	buf := make([]string, 0, n+1)
	for scanner.Scan() {
		buf = append(buf, scanner.Text())
		if len(buf) > n {
			buf = buf[1:]
		}
	}
	return buf
}
