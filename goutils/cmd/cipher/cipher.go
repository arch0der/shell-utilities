// cipher - apply classical ciphers: rot13, rot47, caesar, atbash, vigenere
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func rot13(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z': return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z': return 'a' + (r-'a'+13)%26
		}
		return r
	}, s)
}

func rot47(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '!' && r <= '~' { return '!' + (r-'!'+47)%94 }
		return r
	}, s)
}

func caesar(s string, shift int) string {
	shift = ((shift % 26) + 26) % 26
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z': return 'A' + (r-'A'+rune(shift))%26
		case r >= 'a' && r <= 'z': return 'a' + (r-'a'+rune(shift))%26
		}
		return r
	}, s)
}

func atbash(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z': return 'Z' - (r - 'A')
		case r >= 'a' && r <= 'z': return 'z' - (r - 'a')
		}
		return r
	}, s)
}

func vigenere(s, key string, encode bool) string {
	key = strings.ToLower(key)
	ki := 0
	return strings.Map(func(r rune) rune {
		var base rune
		switch {
		case r >= 'A' && r <= 'Z': base = 'A'
		case r >= 'a' && r <= 'z': base = 'a'
		default: return r
		}
		k := rune(key[ki%len(key)] - 'a')
		ki++
		if encode { return base + (r-base+k)%26 }
		return base + (r-base-k+26)%26
	}, s)
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: cipher <mode> [key] [text...]
  modes: rot13 | rot47 | caesar <n> | atbash | vigenere <key> | devigenere <key>
  If text is omitted, reads from stdin.`)
	os.Exit(1)
}

func apply(fn func(string) string) {
	if len(os.Args) > 3 {
		fmt.Println(fn(strings.Join(os.Args[3:], " ")))
		return
	}
	if len(os.Args) == 3 && os.Args[2] != "" {
		// text may be arg 2 for single-arg modes
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { fmt.Println(fn(sc.Text())) }
}

func applyLine(fn func(string) string) {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { fmt.Println(fn(sc.Text())) }
}

func main() {
	if len(os.Args) < 2 { usage() }
	mode := os.Args[1]

	readText := func(from int) string {
		if len(os.Args) > from { return strings.Join(os.Args[from:], " ") }
		sc := bufio.NewScanner(os.Stdin)
		var lines []string
		for sc.Scan() { lines = append(lines, sc.Text()) }
		return strings.Join(lines, "\n")
	}

	switch mode {
	case "rot13":
		text := readText(2)
		for _, line := range strings.Split(text, "\n") { fmt.Println(rot13(line)) }
	case "rot47":
		text := readText(2)
		for _, line := range strings.Split(text, "\n") { fmt.Println(rot47(line)) }
	case "atbash":
		text := readText(2)
		for _, line := range strings.Split(text, "\n") { fmt.Println(atbash(line)) }
	case "caesar":
		if len(os.Args) < 3 { usage() }
		n, err := strconv.Atoi(os.Args[2])
		if err != nil { fmt.Fprintln(os.Stderr, "caesar: shift must be integer"); os.Exit(1) }
		text := readText(3)
		for _, line := range strings.Split(text, "\n") { fmt.Println(caesar(line, n)) }
	case "vigenere":
		if len(os.Args) < 3 { usage() }
		key := os.Args[2]
		text := readText(3)
		for _, line := range strings.Split(text, "\n") { fmt.Println(vigenere(line, key, true)) }
	case "devigenere":
		if len(os.Args) < 3 { usage() }
		key := os.Args[2]
		text := readText(3)
		for _, line := range strings.Split(text, "\n") { fmt.Println(vigenere(line, key, false)) }
	default:
		// Treat as shorthand: rot13 text...
		_ = applyLine
		fmt.Fprintf(os.Stderr, "cipher: unknown mode %q\n", mode)
		usage()
	}
}
