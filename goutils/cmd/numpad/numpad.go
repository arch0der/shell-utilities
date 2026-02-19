// numpad - pad numbers to fixed width with leading zeros or spaces
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	width := 8
	padChar := "0"
	right := false
	prefix := ""
	suffix := ""

	args := os.Args[1:]
	rest := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-w": i++; width, _ = strconv.Atoi(args[i])
		case "-c": i++; padChar = args[i]
		case "-r": right = true
		case "-p": i++; prefix = args[i]
		case "-s": i++; suffix = args[i]
		default:
			if n, err := strconv.Atoi(args[i]); err == nil { width = n } else { rest = append(rest, args[i]) }
		}
	}

	format := func(s string) string {
		s = strings.TrimSpace(s)
		if s == "" { return s }
		neg := strings.HasPrefix(s, "-")
		if neg { s = s[1:] }
		for len(s) < width { if right { s = s + padChar } else { s = padChar + s } }
		if neg { s = "-" + s }
		return prefix + s + suffix
	}

	if len(rest) > 0 {
		for _, v := range rest { fmt.Println(format(v)) }
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { fmt.Println(format(sc.Text())) }
}
