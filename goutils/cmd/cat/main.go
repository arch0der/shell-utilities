// cat - Concatenate and print files
// Usage: cat [-n] [file...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var numbers = flag.Bool("n", false, "Number all output lines")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cat [-n] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		dump(os.Stdin)
		return
	}
	exitCode := 0
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cat:", err)
			exitCode = 1
			continue
		}
		dump(f)
		f.Close()
	}
	os.Exit(exitCode)
}

var lineNum int

func dump(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineNum++
		if *numbers {
			fmt.Printf("%6d\t%s\n", lineNum, scanner.Text())
		} else {
			fmt.Println(scanner.Text())
		}
	}
}
