// upper - Convert text to UPPERCASE.
//
// Usage:
//
//	upper [FILE...]
//	echo "hello" | upper
//
// Examples:
//
//	echo "hello world" | upper   # HELLO WORLD
//	upper file.txt
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	files := os.Args[1:]
	var readers []*os.File
	if len(files) == 0 {
		readers = []*os.File{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "upper: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for _, r := range readers {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			fmt.Fprintln(w, strings.ToUpper(sc.Text()))
		}
	}
}
