// colorize - highlight pattern matches in stdin with ANSI colours
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var colors = map[string]string{
	"red":     "\033[31m", "green":  "\033[32m", "yellow": "\033[33m",
	"blue":    "\033[34m", "magenta":"\033[35m", "cyan":   "\033[36m",
	"white":   "\033[37m", "bold":   "\033[1m",  "reset":  "\033[0m",
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: colorize <pattern> [color] [file...]")
	fmt.Fprintln(os.Stderr, "  colors: red green yellow blue magenta cyan white bold")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 { usage() }
	pattern := os.Args[1]
	color := "red"
	fileStart := 2
	if len(os.Args) > 2 {
		if _, ok := colors[os.Args[2]]; ok { color = os.Args[2]; fileStart = 3 }
	}
	re, err := regexp.Compile(pattern)
	if err != nil { fmt.Fprintf(os.Stderr, "colorize: invalid pattern: %v\n", err); os.Exit(1) }

	ansi := colors[color]
	reset := colors["reset"]

	highlight := func(line string) string {
		return re.ReplaceAllStringFunc(line, func(m string) string {
			return ansi + m + reset
		})
	}

	process := func(r *os.File) {
		sc := bufio.NewScanner(r)
		for sc.Scan() { fmt.Println(highlight(sc.Text())) }
	}

	if fileStart >= len(os.Args) {
		process(os.Stdin)
		return
	}
	for _, f := range os.Args[fileStart:] {
		_ = strings.TrimSpace(f)
		fh, err := os.Open(f)
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		process(fh); fh.Close()
	}
}
