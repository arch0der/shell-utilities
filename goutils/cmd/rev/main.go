// rev - Reverse characters in each line
// Usage: rev [file...]
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		reverseLines(os.Stdin)
		return
	}
	exitCode := 0
	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "rev:", err)
			exitCode = 1
			continue
		}
		reverseLines(f)
		f.Close()
	}
	os.Exit(exitCode)
}

func reverseLines(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		runes := []rune(scanner.Text())
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		fmt.Println(string(runes))
	}
}
