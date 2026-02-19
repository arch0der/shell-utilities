package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	size := int64(-1)
	noCreate := false
	reference := ""
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-s" && i+1 < len(args):
			i++
			size = parseSize(args[i])
		case strings.HasPrefix(a, "-s"):
			size = parseSize(a[2:])
		case a == "--size" && i+1 < len(args):
			i++
			size = parseSize(args[i])
		case a == "-c" || a == "--no-create":
			noCreate = true
		case a == "-r" && i+1 < len(args):
			i++
			reference = args[i]
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}

	if reference != "" {
		info, err := os.Stat(reference)
		if err != nil {
			fmt.Fprintln(os.Stderr, "truncate:", err)
			os.Exit(1)
		}
		size = info.Size()
	}

	exitCode := 0
	for _, f := range files {
		_, err := os.Stat(f)
		if os.IsNotExist(err) && noCreate {
			continue
		}
		if err := os.Truncate(f, size); err != nil {
			fmt.Fprintf(os.Stderr, "truncate: %s: %v\n", f, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
