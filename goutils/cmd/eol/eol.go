// eol - detect and convert line endings (LF, CRLF, CR)
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func detect(data []byte) string {
	crlf := strings.Count(string(data), "\r\n")
	cr := strings.Count(string(data), "\r") - crlf
	lf := strings.Count(string(data), "\n") - crlf
	switch {
	case crlf > lf && crlf > cr: return "CRLF (Windows)"
	case cr > lf && cr > crlf: return "CR (old Mac)"
	default: return "LF (Unix)"
	}
}

func convert(data []byte, to string) []byte {
	// Normalize to LF first
	s := strings.ReplaceAll(string(data), "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	switch to {
	case "crlf", "windows": s = strings.ReplaceAll(s, "\n", "\r\n")
	case "cr", "mac": s = strings.ReplaceAll(s, "\n", "\r")
	}
	return []byte(s)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: eol <detect|lf|crlf|cr> [file...]")
		os.Exit(1)
	}
	mode := strings.ToLower(os.Args[1])

	readFile := func(path string) ([]byte, error) {
		if path == "-" { return io.ReadAll(os.Stdin) }
		return os.ReadFile(path)
	}
	files := os.Args[2:]
	if len(files) == 0 { files = []string{"-"} }

	for _, f := range files {
		data, err := readFile(f)
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		if mode == "detect" {
			fmt.Printf("%-20s : %s\n", f, detect(data))
		} else {
			out := convert(data, mode)
			if f == "-" { os.Stdout.Write(out) } else {
				if err := os.WriteFile(f, out, 0644); err != nil {
					fmt.Fprintln(os.Stderr, err)
				} else {
					fmt.Fprintf(os.Stderr, "converted %s → %s\n", f, strings.ToUpper(mode))
				}
			}
		}
	}
}
