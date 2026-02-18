// cut - Extract fields or characters from lines
// Usage: cut -f fields [-d delimiter] [file...]
//        cut -c chars [file...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	fields    = flag.String("f", "", "Field numbers to cut (e.g. 1,3 or 1-3)")
	delimiter = flag.String("d", "\t", "Field delimiter")
	chars     = flag.String("c", "", "Character positions to cut (e.g. 1-5)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cut -f fields [-d delim] [file...]\n       cut -c chars [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *fields == "" && *chars == "" {
		flag.Usage()
		os.Exit(1)
	}

	process := func(r io.Reader) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			if *chars != "" {
				fmt.Println(cutChars(line, *chars))
			} else {
				fmt.Println(cutFields(line, *fields, *delimiter))
			}
		}
	}

	if flag.NArg() == 0 {
		process(os.Stdin)
		return
	}
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cut:", err)
			continue
		}
		process(f)
		f.Close()
	}
}

func parseRanges(spec string) [][2]int {
	var ranges [][2]int
	for _, part := range strings.Split(spec, ",") {
		if strings.Contains(part, "-") {
			p := strings.SplitN(part, "-", 2)
			lo, _ := strconv.Atoi(p[0])
			hi, _ := strconv.Atoi(p[1])
			ranges = append(ranges, [2]int{lo, hi})
		} else {
			n, _ := strconv.Atoi(part)
			ranges = append(ranges, [2]int{n, n})
		}
	}
	return ranges
}

func cutFields(line, spec, delim string) string {
	parts := strings.Split(line, delim)
	ranges := parseRanges(spec)
	var out []string
	for i, p := range parts {
		idx := i + 1
		for _, r := range ranges {
			if idx >= r[0] && idx <= r[1] {
				out = append(out, p)
				break
			}
		}
	}
	return strings.Join(out, delim)
}

func cutChars(line, spec string) string {
	ranges := parseRanges(spec)
	runes := []rune(line)
	var out []rune
	for i, c := range runes {
		idx := i + 1
		for _, r := range ranges {
			if idx >= r[0] && idx <= r[1] {
				out = append(out, c)
				break
			}
		}
	}
	return string(out)
}
