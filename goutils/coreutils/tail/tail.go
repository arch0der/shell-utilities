package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	n := 10
	follow := false
	bytes := int64(-1)
	fromStart := false
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-f" || a == "--follow":
			follow = true
		case a == "-n" && i+1 < len(args):
			i++
			s := args[i]
			if strings.HasPrefix(s, "+") {
				fromStart = true
				s = s[1:]
			}
			n, _ = strconv.Atoi(s)
		case strings.HasPrefix(a, "-n"):
			s := a[2:]
			if strings.HasPrefix(s, "+") {
				fromStart = true
				s = s[1:]
			}
			n, _ = strconv.Atoi(s)
		case a == "-c" && i+1 < len(args):
			i++
			bytes = parseSize(args[i])
		case strings.HasPrefix(a, "-c"):
			bytes = parseSize(a[2:])
		case len(a) > 1 && a[0] == '-':
			v, err := strconv.Atoi(a[1:])
			if err == nil {
				n = v
			} else {
				files = append(files, a)
			}
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}

	tailFile := func(r io.Reader, name string, header bool) {
		if header {
			fmt.Printf("==> %s <==\n", name)
		}
		if bytes >= 0 {
			data, _ := io.ReadAll(r)
			start := int64(len(data)) - bytes
			if start < 0 {
				start = 0
			}
			os.Stdout.Write(data[start:])
			return
		}
		var lines []string
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
		if fromStart {
			start := n - 1
			if start < 0 {
				start = 0
			}
			for i := start; i < len(lines); i++ {
				fmt.Println(lines[i])
			}
		} else {
			start := len(lines) - n
			if start < 0 {
				start = 0
			}
			for _, l := range lines[start:] {
				fmt.Println(l)
			}
		}
	}

	if len(files) == 0 {
		tailFile(os.Stdin, "stdin", false)
		return
	}
	for _, f := range files {
		if f == "-" {
			tailFile(os.Stdin, "stdin", len(files) > 1)
			continue
		}
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tail: %s: %v\n", f, err)
			continue
		}
		tailFile(fh, f, len(files) > 1)
		if follow {
			fh.Seek(0, io.SeekEnd)
			for {
				time.Sleep(200 * time.Millisecond)
				sc := bufio.NewScanner(fh)
				for sc.Scan() {
					fmt.Println(sc.Text())
				}
			}
		}
		fh.Close()
	}
}
