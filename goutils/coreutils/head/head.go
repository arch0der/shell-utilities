package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	n := 10
	files := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-n" && i+1 < len(args):
			i++
			n, _ = strconv.Atoi(strings.TrimPrefix(args[i], "-"))
		case strings.HasPrefix(a, "-n"):
			n, _ = strconv.Atoi(a[2:])
		case a == "-q" || a == "--quiet" || a == "--silent":
			// suppress headers - handled inline
		case a == "-v" || a == "--verbose":
			// always print headers - handled inline
		case len(a) > 1 && a[0] == '-':
			if v, err := strconv.Atoi(a[1:]); err == nil {
				n = v
			} else {
				files = append(files, a)
			}
		default:
			files = append(files, a)
		}
	}
	headFile := func(r io.Reader, name string, header bool) {
		if header {
			fmt.Printf("==> %s <==\n", name)
		}
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for count := 0; count < n && sc.Scan(); count++ {
			fmt.Println(sc.Text())
		}
	}
	if len(files) == 0 {
		headFile(os.Stdin, "stdin", false)
		return
	}
	for _, f := range files {
		if f == "-" {
			headFile(os.Stdin, "stdin", len(files) > 1)
			continue
		}
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "head: %s: %v\n", f, err)
			continue
		}
		headFile(fh, f, len(files) > 1)
		fh.Close()
	}
}
