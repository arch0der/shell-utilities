// sed - Stream editor: find and replace
// Usage: sed [-i] 's/pattern/replacement/' [file...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var inPlace = flag.Bool("i", false, "Edit file in-place")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: sed [-i] 's/pattern/replacement/' [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	expr := flag.Arg(0)
	re, repl, global, err := parseExpr(expr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sed:", err)
		os.Exit(1)
	}

	files := flag.Args()[1:]
	if len(files) == 0 {
		process(os.Stdin, os.Stdout, re, repl, global)
		return
	}

	for _, path := range files {
		if *inPlace {
			processInPlace(path, re, repl, global)
		} else {
			f, err := os.Open(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "sed:", err)
				continue
			}
			process(f, os.Stdout, re, repl, global)
			f.Close()
		}
	}
}

func parseExpr(expr string) (*regexp.Regexp, string, bool, error) {
	if len(expr) < 2 {
		return nil, "", false, fmt.Errorf("invalid expression")
	}
	sep := string(expr[1])
	parts := strings.Split(expr[2:], sep)
	if len(parts) < 2 {
		return nil, "", false, fmt.Errorf("invalid s expression")
	}
	pattern, repl := parts[0], parts[1]
	flags := ""
	if len(parts) > 2 {
		flags = parts[2]
	}
	global := strings.Contains(flags, "g")
	re, err := regexp.Compile(pattern)
	return re, repl, global, err
}

func process(in io.Reader, out io.Writer, re *regexp.Regexp, repl string, global bool) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if global {
			line = re.ReplaceAllString(line, repl)
		} else {
			loc := re.FindStringIndex(line)
			if loc != nil {
				line = line[:loc[0]] + re.ReplaceAllString(line[loc[0]:loc[1]], repl) + line[loc[1]:]
			}
		}
		fmt.Fprintln(out, line)
	}
}

func processInPlace(path string, re *regexp.Regexp, repl string, global bool) {
	in, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sed:", err)
		return
	}
	tmp, err := os.CreateTemp("", "sed")
	if err != nil {
		in.Close()
		return
	}
	process(in, tmp, re, repl, global)
	in.Close()
	tmp.Close()
	os.Rename(tmp.Name(), path)
}
