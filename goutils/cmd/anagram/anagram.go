// anagram - check if two words are anagrams, or group anagram sets from stdin
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func key(s string) string {
	r := []rune(strings.ToLower(s))
	sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
	return string(r)
}

func main() {
	if len(os.Args) == 3 {
		a, b := os.Args[1], os.Args[2]
		if key(a) == key(b) {
			fmt.Printf("✓ %q and %q ARE anagrams\n", a, b)
		} else {
			fmt.Printf("✗ %q and %q are NOT anagrams\n", a, b)
			os.Exit(1)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] != "-" {
		fmt.Fprintln(os.Stderr, "usage: anagram <word1> <word2>  |  anagram < wordlist")
		os.Exit(2)
	}
	groups := map[string][]string{}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		w := strings.TrimSpace(sc.Text())
		if w != "" {
			k := key(w)
			groups[k] = append(groups[k], w)
		}
	}
	keys := make([]string, 0, len(groups))
	for k := range groups { keys = append(keys, k) }
	sort.Strings(keys)
	found := false
	for _, k := range keys {
		if len(groups[k]) > 1 {
			fmt.Println(strings.Join(groups[k], "  "))
			found = true
		}
	}
	if !found {
		fmt.Println("No anagram groups found.")
	}
}
