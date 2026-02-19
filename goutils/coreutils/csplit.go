package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func init() { register("csplit", runCsplit) }

func runCsplit() {
	args := os.Args[1:]
	prefix := "xx"
	silent := false
	keepFiles := false
	files := []string{}
	patterns := []string{}

	i := 0
	for i < len(args) {
		a := args[i]
		switch {
		case a == "-f" && i+1 < len(args):
			i++
			prefix = args[i]
		case strings.HasPrefix(a, "-f"):
			prefix = a[2:]
		case a == "-s" || a == "--silent" || a == "--quiet":
			silent = true
		case a == "-k" || a == "--keep-files":
			keepFiles = true
		case !strings.HasPrefix(a, "-"):
			if len(files) == 0 {
				files = append(files, a)
			} else {
				patterns = append(patterns, a)
			}
		}
		i++
	}
	_ = keepFiles

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "csplit: missing operand")
		os.Exit(1)
	}

	var r *os.File
	if files[0] == "-" {
		r = os.Stdin
	} else {
		var err error
		r, err = os.Open(files[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "csplit: %v\n", err)
			os.Exit(1)
		}
		defer r.Close()
	}

	var lines []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	// Split points: line numbers or regex
	splits := []int{0}
	for _, p := range patterns {
		if strings.HasPrefix(p, "/") && strings.HasSuffix(p, "/") {
			re := regexp.MustCompile(p[1 : len(p)-1])
			for idx, l := range lines {
				if re.MatchString(l) {
					splits = append(splits, idx)
					break
				}
			}
		} else if n, err := strconv.Atoi(p); err == nil {
			splits = append(splits, n-1)
		}
	}
	splits = append(splits, len(lines))

	fileNum := 0
	for i := 0; i < len(splits)-1; i++ {
		start, end := splits[i], splits[i+1]
		if end > len(lines) {
			end = len(lines)
		}
		chunk := strings.Join(lines[start:end], "\n")
		if end > start {
			chunk += "\n"
		}
		name := fmt.Sprintf("%s%02d", prefix, fileNum)
		if err := os.WriteFile(name, []byte(chunk), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "csplit: %v\n", err)
		}
		if !silent {
			fmt.Println(len(chunk))
		}
		fileNum++
	}
}
