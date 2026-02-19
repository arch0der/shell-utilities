// fmtgo - format Go source code and report formatting differences
package main

import (
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"
)

func main() {
	diff := false
	write := false
	list := false
	var files []string

	for _, arg := range os.Args[1:] {
		switch arg {
		case "-d", "--diff": diff = true
		case "-w", "--write": write = true
		case "-l", "--list": list = true
		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Fprintln(os.Stderr, "usage: fmtgo [-d|-w|-l] [file...]")
				os.Exit(1)
			}
			files = append(files, arg)
		}
	}

	process := func(src []byte, path string) {
		formatted, err := format.Source(src)
		if err != nil { fmt.Fprintf(os.Stderr, "fmtgo: %s: %v\n", path, err); return }
		if string(formatted) == string(src) {
			if list && path != "<stdin>" { fmt.Println(path + " (unchanged)") }
			return
		}
		if list { fmt.Println(path); return }
		if diff {
			showDiff(string(src), string(formatted), path)
			return
		}
		if write && path != "<stdin>" {
			os.WriteFile(path, formatted, 0644)
			fmt.Fprintln(os.Stderr, "formatted:", path)
			return
		}
		os.Stdout.Write(formatted)
	}

	if len(files) == 0 {
		src, _ := io.ReadAll(os.Stdin); process(src, "<stdin>"); return
	}
	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		process(src, f)
	}
}

func showDiff(a, b, name string) {
	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")
	fmt.Printf("--- %s (original)\n+++ %s (formatted)\n", name, name)
	for i := 0; i < len(aLines) || i < len(bLines); i++ {
		al, bl := "", ""
		if i < len(aLines) { al = aLines[i] }
		if i < len(bLines) { bl = bLines[i] }
		if al != bl {
			if al != "" { fmt.Printf("-%s\n", al) }
			if bl != "" { fmt.Printf("+%s\n", bl) }
		}
	}
}
