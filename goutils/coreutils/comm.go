package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func init() { register("comm", runComm) }

func runComm() {
	args := os.Args[1:]
	suppress := map[int]bool{}
	files := []string{}
	for _, a := range args {
		switch a {
		case "-1":
			suppress[1] = true
		case "-2":
			suppress[2] = true
		case "-3":
			suppress[3] = true
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}
	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "comm: missing operand")
		os.Exit(1)
	}
	readLines := func(f string) []string {
		var r *os.File
		if f == "-" {
			r = os.Stdin
		} else {
			var err error
			r, err = os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "comm: %v\n", err)
				os.Exit(1)
			}
			defer r.Close()
		}
		var lines []string
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
		return lines
	}
	l1 := readLines(files[0])
	l2 := readLines(files[1])
	i, j := 0, 0
	for i < len(l1) || j < len(l2) {
		var cmp int
		if i >= len(l1) {
			cmp = 1
		} else if j >= len(l2) {
			cmp = -1
		} else {
			cmp = strings.Compare(l1[i], l2[j])
		}
		if cmp < 0 {
			if !suppress[1] {
				fmt.Println(l1[i])
			}
			i++
		} else if cmp > 0 {
			if !suppress[2] {
				fmt.Printf("\t%s\n", l2[j])
			}
			j++
		} else {
			if !suppress[3] {
				fmt.Printf("\t\t%s\n", l1[i])
			}
			i++
			j++
		}
	}
}
