package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func init() { register("uniq", runUniq) }

func runUniq() {
	args := os.Args[1:]
	countMode := false
	uniqueOnly := false
	dupOnly := false
	ignoreCase := false
	skipFields := 0
	skipChars := 0
	zero := false
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-c" || a == "--count":
			countMode = true
		case a == "-u" || a == "--unique":
			uniqueOnly = true
		case a == "-d" || a == "--repeated":
			dupOnly = true
		case a == "-i" || a == "--ignore-case":
			ignoreCase = true
		case a == "-z" || a == "--zero-terminated":
			zero = true
		case a == "-f" && i+1 < len(args):
			i++
			fmt.Sscan(args[i], &skipFields)
		case strings.HasPrefix(a, "-f"):
			fmt.Sscan(a[2:], &skipFields)
		case a == "-s" && i+1 < len(args):
			i++
			fmt.Sscan(args[i], &skipChars)
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}
	_ = skipFields
	_ = skipChars
	_ = zero

	var input, output *os.File = os.Stdin, os.Stdout
	if len(files) > 0 {
		fh, err := os.Open(files[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "uniq:", err)
			os.Exit(1)
		}
		defer fh.Close()
		input = fh
	}
	if len(files) > 1 {
		fh, err := os.Create(files[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "uniq:", err)
			os.Exit(1)
		}
		defer fh.Close()
		output = fh
	}

	type group struct {
		line  string
		count int
	}
	var groups []group
	sc := bufio.NewScanner(input)
	for sc.Scan() {
		line := sc.Text()
		key := line
		if ignoreCase {
			key = strings.ToLower(key)
		}
		if len(groups) > 0 {
			prevKey := groups[len(groups)-1].line
			if ignoreCase {
				prevKey = strings.ToLower(prevKey)
			}
			if prevKey == key {
				groups[len(groups)-1].count++
				continue
			}
		}
		groups = append(groups, group{line, 1})
	}

	bw := bufio.NewWriter(output)
	defer bw.Flush()

	for _, g := range groups {
		show := true
		if uniqueOnly && g.count > 1 {
			show = false
		}
		if dupOnly && g.count == 1 {
			show = false
		}
		if !show {
			continue
		}
		if countMode {
			fmt.Fprintf(bw, "%7d %s\n", g.count, g.line)
		} else {
			io.WriteString(bw, g.line+"\n")
		}
	}
}
