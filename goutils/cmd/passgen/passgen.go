// passgen - generate secure passwords with configurable rules
package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

const (
	lower   = "abcdefghijklmnopqrstuvwxyz"
	upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
	symbols = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

func randChar(charset string) (byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil { return 0, err }
	return charset[n.Int64()], nil
}

func generate(length int, useLower, useUpper, useDigits, useSymbols bool, exclude string) (string, error) {
	var charset strings.Builder
	var mustHave []string
	if useLower { charset.WriteString(lower); mustHave = append(mustHave, lower) }
	if useUpper { charset.WriteString(upper); mustHave = append(mustHave, upper) }
	if useDigits { charset.WriteString(digits); mustHave = append(mustHave, digits) }
	if useSymbols { charset.WriteString(symbols); mustHave = append(mustHave, symbols) }
	if charset.Len() == 0 { return "", fmt.Errorf("no character classes selected") }

	// Remove excluded chars
	cs := charset.String()
	for _, ch := range exclude { cs = strings.ReplaceAll(cs, string(ch), "") }
	if len(cs) == 0 { return "", fmt.Errorf("charset empty after exclusions") }

	for {
		buf := make([]byte, length)
		for i := range buf {
			ch, err := randChar(cs)
			if err != nil { return "", err }
			buf[i] = ch
		}
		// Check all required classes present
		pw := string(buf)
		ok := true
		for _, must := range mustHave {
			found := false
			for _, ch := range must {
				if strings.ContainsRune(pw, ch) { found = true; break }
			}
			if !found { ok = false; break }
		}
		if ok { return pw, nil }
	}
}

func main() {
	length := 20
	count := 1
	useLower, useUpper, useDigits, useSymbols := true, true, true, true
	exclude := ""

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-l": i++; length, _ = strconv.Atoi(args[i])
		case "-n": i++; count, _ = strconv.Atoi(args[i])
		case "-x": i++; exclude = args[i]
		case "--no-lower": useLower = false
		case "--no-upper": useUpper = false
		case "--no-digits": useDigits = false
		case "--no-symbols": useSymbols = false
		case "--only-alpha": useDigits = false; useSymbols = false
		case "--only-alnum": useSymbols = false
		default:
			if n, err := strconv.Atoi(args[i]); err == nil { length = n }
		}
	}

	for i := 0; i < count; i++ {
		pw, err := generate(length, useLower, useUpper, useDigits, useSymbols, exclude)
		if err != nil { fmt.Fprintln(os.Stderr, "passgen:", err); os.Exit(1) }
		fmt.Println(pw)
	}
}
