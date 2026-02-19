// linemap - transform each line with sed-like operations: map, split, join, swap fields
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: linemap <operation> [args] [file...]
  Operations:
    prefix <str>           prepend string to each line
    suffix <str>           append string to each line
    replace <old> <new>    replace literal string
    sub <pat> <rep>        regex substitution
    field <n> [delim]      extract field n (1-based, default delim: tab)
    swap <n> <m> [delim]   swap fields n and m
    upper / lower          change case
    reverse                reverse each line
    len                    print length of each line
    count                  count lines (print at end)
    num                    prefix line number`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 { usage() }
	op := os.Args[1]
	args := os.Args[2:]

	var transform func(string, int) string
	finalAction := func() {}

	switch op {
	case "prefix":
		if len(args) < 1 { usage() }
		p := args[0]; transform = func(line string, _ int) string { return p + line }
	case "suffix":
		if len(args) < 1 { usage() }
		s := args[0]; transform = func(line string, _ int) string { return line + s }
	case "replace":
		if len(args) < 2 { usage() }
		transform = func(line string, _ int) string { return strings.ReplaceAll(line, args[0], args[1]) }
	case "sub":
		if len(args) < 2 { usage() }
		re, err := regexp.Compile(args[0]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		transform = func(line string, _ int) string { return re.ReplaceAllString(line, args[1]) }
	case "field":
		if len(args) < 1 { usage() }
		n, _ := strconv.Atoi(args[0]); delim := "\t"
		if len(args) > 1 { delim = args[1] }
		transform = func(line string, _ int) string {
			parts := strings.Split(line, delim)
			if n < 1 || n > len(parts) { return "" }; return parts[n-1]
		}
	case "swap":
		if len(args) < 2 { usage() }
		a, _ := strconv.Atoi(args[0]); b, _ := strconv.Atoi(args[1]); delim := "\t"
		if len(args) > 2 { delim = args[2] }
		transform = func(line string, _ int) string {
			parts := strings.Split(line, delim)
			if a < 1 || a > len(parts) || b < 1 || b > len(parts) { return line }
			parts[a-1], parts[b-1] = parts[b-1], parts[a-1]
			return strings.Join(parts, delim)
		}
	case "upper": transform = func(line string, _ int) string { return strings.ToUpper(line) }
	case "lower": transform = func(line string, _ int) string { return strings.ToLower(line) }
	case "reverse":
		transform = func(line string, _ int) string {
			runes := []rune(line)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 { runes[i], runes[j] = runes[j], runes[i] }
			return string(runes)
		}
	case "len": transform = func(line string, _ int) string { return strconv.Itoa(len([]rune(line))) }
	case "num": transform = func(line string, n int) string { return fmt.Sprintf("%6d  %s", n, line) }
	case "count":
		count := 0
		transform = func(line string, _ int) string { count++; return "" }
		finalAction = func() { fmt.Println(count) }
	default:
		fmt.Fprintf(os.Stderr, "linemap: unknown operation %q\n", op); usage()
	}

	process := func(r *os.File) {
		sc := bufio.NewScanner(r); n := 0
		for sc.Scan() {
			n++; result := transform(sc.Text(), n)
			if op != "count" { fmt.Println(result) }
		}
	}

	files := args
	if op == "prefix" || op == "suffix" || op == "replace" || op == "sub" { files = args[1:] }
	if op == "field" { files = args[1:] }
	if op == "swap" { files = args[2:] }

	if len(files) == 0 { process(os.Stdin) } else {
		for _, f := range files {
			fh, err := os.Open(f); if err != nil { fmt.Fprintln(os.Stderr, err); continue }
			process(fh); fh.Close()
		}
	}
	finalAction()
}
