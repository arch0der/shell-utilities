package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	field1, field2 := 1, 1
	delim := " "
	empty := ""
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-1" && i+1 < len(args):
			i++
			fmt.Sscan(args[i], &field1)
		case a == "-2" && i+1 < len(args):
			i++
			fmt.Sscan(args[i], &field2)
		case (a == "-t" || a == "-d") && i+1 < len(args):
			i++
			delim = args[i]
		case a == "-e" && i+1 < len(args):
			i++
			empty = args[i]
		case a == "-a":
			i++ // skip file number arg
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}

	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "join: missing operand")
		os.Exit(1)
	}

	readFile := func(f string, field int) map[string][]string {
		m := map[string][]string{}
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "join: %v\n", err)
			os.Exit(1)
		}
		defer fh.Close()
		sc := bufio.NewScanner(fh)
		for sc.Scan() {
			line := sc.Text()
			var parts []string
			if delim == " " {
				parts = strings.Fields(line)
			} else {
				parts = strings.Split(line, delim)
			}
			if len(parts) < field {
				continue
			}
			key := parts[field-1]
			m[key] = parts
		}
		return m
	}

	m1 := readFile(files[0], field1)
	m2 := readFile(files[1], field2)

	seen := map[string]bool{}
	for key, parts1 := range m1 {
		seen[key] = true
		parts2, ok := m2[key]
		if !ok {
			parts2 = []string{empty}
		}
		var out []string
		out = append(out, key)
		for i, p := range parts1 {
			if i+1 != field1 {
				out = append(out, p)
			}
		}
		for i, p := range parts2 {
			if i+1 != field2 {
				out = append(out, p)
			}
		}
		fmt.Println(strings.Join(out, delim))
	}
	_ = seen
}
