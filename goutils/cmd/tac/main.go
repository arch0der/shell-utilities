// tac - Print file lines in reverse order
// Usage: tac [file...]
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printReverse(os.Stdin)
		return
	}
	exitCode := 0
	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tac:", err)
			exitCode = 1
			continue
		}
		printReverse(f)
		f.Close()
	}
	os.Exit(exitCode)
}

func printReverse(r io.Reader) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for i := len(lines) - 1; i >= 0; i-- {
		fmt.Println(lines[i])
	}
}
