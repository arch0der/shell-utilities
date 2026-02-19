// comm - Compare two sorted files line by line
// Usage: comm [-123] file1 file2
// Output: col1=only-in-file1, col2=only-in-file2, col3=in-both
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var (
	suppress1 = flag.Bool("1", false, "Suppress lines only in file1")
	suppress2 = flag.Bool("2", false, "Suppress lines only in file2")
	suppress3 = flag.Bool("3", false, "Suppress lines in both files")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: comm [-123] file1 file2")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	lines1 := readLines(flag.Arg(0))
	lines2 := readLines(flag.Arg(1))

	i, j := 0, 0
	for i < len(lines1) || j < len(lines2) {
		var cmp int
		if i >= len(lines1) {
			cmp = 1
		} else if j >= len(lines2) {
			cmp = -1
		} else if lines1[i] < lines2[j] {
			cmp = -1
		} else if lines1[i] > lines2[j] {
			cmp = 1
		} else {
			cmp = 0
		}

		switch cmp {
		case -1:
			if !*suppress1 {
				fmt.Println(lines1[i])
			}
			i++
		case 1:
			if !*suppress2 {
				fmt.Println("\t" + lines2[j])
			}
			j++
		case 0:
			if !*suppress3 {
				fmt.Println("\t\t" + lines1[i])
			}
			i++
			j++
		}
	}
}

func readLines(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "comm:", err)
		os.Exit(1)
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
