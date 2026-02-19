// zigzag - encode/decode text using the Rail Fence (ZigZag) cipher
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func encode(text string, rails int) string {
	if rails < 2 { return text }
	fence := make([]strings.Builder, rails)
	rail, dir := 0, 1
	for _, ch := range text {
		fence[rail].WriteRune(ch)
		if rail == 0 { dir = 1 } else if rail == rails-1 { dir = -1 }
		rail += dir
	}
	var result strings.Builder
	for _, f := range fence { result.WriteString(f.String()) }
	return result.String()
}

func decode(cipher string, rails int) string {
	if rails < 2 { return cipher }
	n := len(cipher)
	// Determine pattern
	pattern := make([]int, n)
	rail, dir := 0, 1
	for i := range pattern {
		pattern[i] = rail
		if rail == 0 { dir = 1 } else if rail == rails-1 { dir = -1 }
		rail += dir
	}
	// Count chars per rail
	counts := make([]int, rails)
	for _, r := range pattern { counts[r]++ }
	// Slice cipher into rails
	railStrings := make([][]rune, rails)
	pos := 0
	for r := 0; r < rails; r++ {
		railStrings[r] = []rune(cipher[pos : pos+counts[r]])
		pos += counts[r]
	}
	// Read in zigzag order
	railIdx := make([]int, rails)
	result := make([]rune, n)
	for i, r := range pattern {
		result[i] = railStrings[r][railIdx[r]]
		railIdx[r]++
	}
	return string(result)
}

func main() {
	decode_ := false
	rails := 3
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d", "decode": decode_ = true
		case "-r", "--rails": i++; rails, _ = strconv.Atoi(args[i])
		default:
			if n, err := strconv.Atoi(args[i]); err == nil { rails = n }
		}
	}

	process := func(s string) string {
		if decode_ { return decode(s, rails) }
		return encode(s, rails)
	}

	textArgs := []string{}
	for _, a := range args {
		if a == "-d" || a == "decode" || a == "-r" || a == "--rails" { continue }
		if _, err := strconv.Atoi(a); err == nil { continue }
		textArgs = append(textArgs, a)
	}

	if len(textArgs) > 0 {
		fmt.Println(process(strings.Join(textArgs, " ")))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { fmt.Println(process(sc.Text())) }
}
