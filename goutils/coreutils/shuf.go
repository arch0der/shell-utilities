package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() { register("shuf", runShuf) }

func runShuf() {
	args := os.Args[1:]
	count := -1
	repeat := false
	zero := false
	echo := false
	inputRange := ""
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-n" && i+1 < len(args):
			i++
			count, _ = strconv.Atoi(args[i])
		case strings.HasPrefix(a, "-n"):
			count, _ = strconv.Atoi(a[2:])
		case a == "-r" || a == "--repeat":
			repeat = true
		case a == "-z" || a == "--zero-terminated":
			zero = true
		case a == "-e" || a == "--echo":
			echo = true
		case a == "-i" && i+1 < len(args):
			i++
			inputRange = args[i]
		case strings.HasPrefix(a, "-i"):
			inputRange = a[2:]
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sep := "\n"
	if zero {
		sep = "\x00"
	}

	var lines []string

	if inputRange != "" {
		parts := strings.SplitN(inputRange, "-", 2)
		lo, _ := strconv.Atoi(parts[0])
		hi, _ := strconv.Atoi(parts[1])
		for i := lo; i <= hi; i++ {
			lines = append(lines, strconv.Itoa(i))
		}
	} else if echo {
		lines = files
	} else {
		readLines := func(rd io.Reader) {
			sc := bufio.NewScanner(rd)
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
			} else {
				fh, _ := os.Open(f)
				readLines(fh)
				fh.Close()
			}
		}
	}

	if repeat {
		if count < 0 {
			count = len(lines)
		}
		for i := 0; i < count; i++ {
			fmt.Print(lines[r.Intn(len(lines))] + sep)
		}
		return
	}

	// Shuffle
	r.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })

	if count >= 0 && count < len(lines) {
		lines = lines[:count]
	}
	for i, l := range lines {
		if i < len(lines)-1 || zero {
			fmt.Print(l + sep)
		} else {
			fmt.Println(l)
		}
	}
}
