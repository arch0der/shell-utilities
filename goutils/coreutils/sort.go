package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func init() { register("sort", runSort) }

func runSort() {
	args := os.Args[1:]
	reverse := false
	unique := false
	numeric := false
	ignoreCase := false
	monthSort := false
	humanSort := false
	versionSort := false
	randomSort := false
	stable := false
	checkMode := false
	output := ""
	key := ""
	fieldSep := ""
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-r" || a == "--reverse":
			reverse = true
		case a == "-u" || a == "--unique":
			unique = true
		case a == "-n" || a == "--numeric-sort":
			numeric = true
		case a == "-f" || a == "--ignore-case":
			ignoreCase = true
		case a == "-M" || a == "--month-sort":
			monthSort = true
		case a == "-h" || a == "--human-numeric-sort":
			humanSort = true
		case a == "-V" || a == "--version-sort":
			versionSort = true
		case a == "-R" || a == "--random-sort":
			randomSort = true
		case a == "-s" || a == "--stable":
			stable = true
		case a == "-c" || a == "--check":
			checkMode = true
		case a == "-o" && i+1 < len(args):
			i++
			output = args[i]
		case strings.HasPrefix(a, "-o"):
			output = a[2:]
		case a == "-k" && i+1 < len(args):
			i++
			key = args[i]
		case strings.HasPrefix(a, "-k"):
			key = a[2:]
		case a == "-t" && i+1 < len(args):
			i++
			fieldSep = args[i]
		case strings.HasPrefix(a, "-t"):
			fieldSep = a[2:]
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}
	_ = monthSort
	_ = humanSort
	_ = versionSort
	_ = randomSort
	_ = stable
	_ = key
	_ = fieldSep

	var lines []string
	readLines := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
	}

	if len(files) == 0 {
		readLines(os.Stdin)
	}
	for _, f := range files {
		if f == "-" {
			readLines(os.Stdin)
			continue
		}
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sort: %s: %v\n", f, err)
			continue
		}
		readLines(fh)
		fh.Close()
	}

	if checkMode {
		for i := 1; i < len(lines); i++ {
			a, b := lines[i-1], lines[i]
			if reverse {
				a, b = b, a
			}
			if strings.Compare(a, b) > 0 {
				fmt.Fprintf(os.Stderr, "sort: %s:%d: disorder: %s\n", "", i+1, lines[i])
				os.Exit(1)
			}
		}
		return
	}

	lessFunc := func(a, b string) bool {
		if ignoreCase {
			a, b = strings.ToLower(a), strings.ToLower(b)
		}
		if numeric {
			an, aerr := strconv.ParseFloat(strings.Fields(a+"\t")[0], 64)
			bn, berr := strconv.ParseFloat(strings.Fields(b+"\t")[0], 64)
			if aerr == nil && berr == nil {
				return an < bn
			}
		}
		return a < b
	}

	sort.SliceStable(lines, func(i, j int) bool {
		l := lessFunc(lines[i], lines[j])
		if reverse {
			return !l && lines[i] != lines[j]
		}
		return l
	})

	if unique {
		var out []string
		seen := map[string]bool{}
		for _, l := range lines {
			key2 := l
			if ignoreCase {
				key2 = strings.ToLower(l)
			}
			if !seen[key2] {
				out = append(out, l)
				seen[key2] = true
			}
		}
		lines = out
	}

	var w io.Writer = os.Stdout
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			fmt.Fprintln(os.Stderr, "sort:", err)
			os.Exit(1)
		}
		defer f.Close()
		w = f
	}

	bw := bufio.NewWriter(w)
	for _, l := range lines {
		bw.WriteString(l + "\n")
	}
	bw.Flush()
}
