// palindrome - check if text is a palindrome, or find palindromes in stdin
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func normalize(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) { b.WriteRune(r) }
	}
	return b.String()
}

func isPalindrome(s string) bool {
	n := normalize(s)
	runes := []rune(n)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] { return false }
	}
	return len(runes) > 0
}

func main() {
	if len(os.Args) > 1 {
		text := strings.Join(os.Args[1:], " ")
		if isPalindrome(text) {
			fmt.Printf("✓ %q is a palindrome\n", text)
		} else {
			fmt.Printf("✗ %q is NOT a palindrome\n", text)
			os.Exit(1)
		}
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		if isPalindrome(line) { fmt.Println(line) }
	}
}
