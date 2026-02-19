// linesplit - split long lines at word boundaries or hard column limit
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func wrap(line string, width int, hard bool, indent string) []string {
	if utf8.RuneCountInString(line) <= width { return []string{line} }
	if hard {
		var out []string
		runes := []rune(line)
		for len(runes) > width { out = append(out, string(runes[:width])); runes = runes[width:] }
		if len(runes) > 0 { out = append(out, string(runes)) }
		return out
	}
	words := strings.Fields(line)
	if len(words) == 0 { return []string{""} }
	var lines []string
	cur := ""
	for _, w := range words {
		if cur == "" {
			cur = indent + w
		} else if utf8.RuneCountInString(cur)+1+utf8.RuneCountInString(w) > width {
			lines = append(lines, cur)
			cur = indent + w
		} else {
			cur += " " + w
		}
	}
	if cur != "" { lines = append(lines, cur) }
	return lines
}

func main() {
	width := 80
	hard := false
	indent := ""
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-w": i++; width, _ = strconv.Atoi(args[i])
		case "-H", "--hard": hard = true
		case "--indent": i++; indent = args[i]
		default:
			if n, err := strconv.Atoi(args[i]); err == nil { width = n }
		}
	}

	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 10*1024*1024), 10*1024*1024)
	for sc.Scan() {
		for _, l := range wrap(sc.Text(), width, hard, indent) { fmt.Println(l) }
	}
}
