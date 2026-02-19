// chunk - split stdin into chunks of N lines, separated by a delimiter
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: chunk <n> [separator]")
		os.Exit(1)
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n < 1 {
		fmt.Fprintln(os.Stderr, "chunk: n must be a positive integer"); os.Exit(1)
	}
	sep := "---"
	if len(os.Args) >= 3 { sep = os.Args[2] }

	sc := bufio.NewScanner(os.Stdin)
	count, first := 0, true
	for sc.Scan() {
		if count == 0 && !first { fmt.Println(sep) }
		first = false
		fmt.Println(sc.Text())
		count++
		if count >= n { count = 0 }
	}
}
