package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	multiple := false
	zero := false
	suffix := ""
	paths := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-a" || a == "--multiple":
			multiple = true
		case a == "-z" || a == "--zero":
			zero = true
		case a == "-s" && i+1 < len(args):
			i++
			suffix = args[i]
			multiple = true
		case strings.HasPrefix(a, "-s"):
			suffix = a[2:]
			multiple = true
		case !strings.HasPrefix(a, "-"):
			paths = append(paths, a)
		}
	}
	sep := "\n"
	if zero {
		sep = "\x00"
	}
	process := func(p string) string {
		b := filepath.Base(p)
		if suffix != "" {
			b = strings.TrimSuffix(b, suffix)
		}
		return b
	}
	if !multiple && len(paths) >= 2 {
		suffix = paths[1]
		paths = paths[:1]
	}
	for i, p := range paths {
		result := process(p)
		if i < len(paths)-1 {
			fmt.Print(result + sep)
		} else {
			if zero {
				fmt.Print(result + sep)
			} else {
				fmt.Println(result)
			}
		}
	}
}
