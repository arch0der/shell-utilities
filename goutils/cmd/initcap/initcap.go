// initcap - capitalize first letter of each word (title case), respecting common words
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Words that should stay lowercase (unless first or last word)
var lowercase = map[string]bool{
	"a":true,"an":true,"the":true,"and":true,"but":true,"or":true,"nor":true,
	"for":true,"yet":true,"so":true,"at":true,"by":true,"for":true,"in":true,
	"of":true,"on":true,"to":true,"up":true,"as":true,"vs":true,"via":true,
}

func titleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		lower := strings.ToLower(w)
		if i == 0 || i == len(words)-1 || !lowercase[lower] {
			runes := []rune(w)
			if len(runes) > 0 { runes[0] = unicode.ToUpper(runes[0]) }
			words[i] = string(runes)
		} else {
			words[i] = lower
		}
	}
	return strings.Join(words, " ")
}

func main() {
	if len(os.Args) > 1 {
		fmt.Println(titleCase(strings.Join(os.Args[1:], " ")))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { fmt.Println(titleCase(sc.Text())) }
}
