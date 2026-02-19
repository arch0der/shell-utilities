// isutf8 - check if files or stdin are valid UTF-8, report invalid byte positions
package main

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

func check(r io.Reader, name string) bool {
	data, err := io.ReadAll(r)
	if err != nil { fmt.Fprintln(os.Stderr, err); return false }
	if utf8.Valid(data) {
		fmt.Printf("✓ %s  (valid UTF-8, %d bytes)\n", name, len(data))
		return true
	}
	// Find offending bytes
	fmt.Printf("✗ %s  (INVALID UTF-8, %d bytes)\n", name, len(data))
	line, col, byteOff := 1, 1, 0
	for i := 0; i < len(data); {
		r, size := utf8.DecodeRune(data[i:])
		if r == utf8.RuneError && size == 1 {
			fmt.Printf("  offset:%d  line:%d  col:%d  byte:0x%02X\n", byteOff+i, line, col, data[i])
		}
		if data[i] == '\n' { line++; col = 1 } else { col++ }
		i += size; byteOff = 0
	}
	return false
}

func main() {
	ok := true
	if len(os.Args) == 1 {
		if !check(os.Stdin, "<stdin>") { ok = false }
	} else {
		for _, path := range os.Args[1:] {
			f, err := os.Open(path)
			if err != nil { fmt.Fprintln(os.Stderr, err); ok = false; continue }
			if !check(f, path) { ok = false }
			f.Close()
		}
	}
	if !ok { os.Exit(1) }
}
