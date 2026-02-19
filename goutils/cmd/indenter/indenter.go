// indenter - add consistent indentation to text blocks
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	indent := "  "
	level := 1
	prefix := ""
	strip := false
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-t", "--tab": indent = "\t"
		case "-s", "--strip": strip = true
		case "-p": i++; prefix = args[i]
		case "-l", "--level": i++; level, _ = strconv.Atoi(args[i])
		default:
			if n, err := strconv.Atoi(args[i]); err == nil { level = n } else {
				indent = args[i]
			}
		}
	}
	pad := strings.Repeat(indent, level) + prefix
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		if strip {
			// Remove up to len(pad) whitespace chars from start
			trimmed := strings.TrimLeft(line, " \t")
			removed := len(line) - len(trimmed)
			if removed > len(pad) { removed = len(pad) }
			line = line[removed:]
		} else {
			if line != "" { line = pad + line }
		}
		fmt.Println(line)
	}
}
