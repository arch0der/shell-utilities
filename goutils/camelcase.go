// camelcase - Convert strings to camelCase or PascalCase.
//
// Usage:
//
//	camelcase [OPTIONS] [TEXT...]
//	echo "hello_world" | camelcase
//
// Options:
//
//	-p        PascalCase (UpperCamelCase)
//	-s SEP    Input word separator (default: auto-detect _-. and spaces)
//
// Examples:
//
//	echo "hello_world"   | camelcase      # helloWorld
//	echo "hello_world"   | camelcase -p   # HelloWorld
//	echo "my-css-class"  | camelcase      # myCssClass
//	camelcase "the quick brown fox"        # theQuickBrownFox
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var (
	pascal = flag.Bool("p", false, "PascalCase")
	sep    = flag.String("s", "", "separator")
)

func toCamel(s string, upper bool) string {
	var words []string
	if *sep != "" {
		words = strings.Split(s, *sep)
	} else {
		// auto split on _, -, ., space, and case transitions
		var cur strings.Builder
		runes := []rune(s)
		for i, r := range runes {
			if r == '_' || r == '-' || r == '.' || r == ' ' {
				if cur.Len() > 0 {
					words = append(words, cur.String())
					cur.Reset()
				}
			} else if i > 0 && unicode.IsUpper(r) && unicode.IsLower(runes[i-1]) {
				words = append(words, cur.String())
				cur.Reset()
				cur.WriteRune(r)
			} else {
				cur.WriteRune(r)
			}
		}
		if cur.Len() > 0 {
			words = append(words, cur.String())
		}
	}

	var result strings.Builder
	for i, w := range words {
		if w == "" {
			continue
		}
		w = strings.ToLower(w)
		if i == 0 && !upper {
			result.WriteString(w)
		} else {
			result.WriteString(strings.ToUpper(w[:1]) + w[1:])
		}
	}
	return result.String()
}

func main() {
	flag.Parse()
	args := flag.Args()
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if len(args) > 0 {
		fmt.Fprintln(w, toCamel(strings.Join(args, " "), *pascal))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Fprintln(w, toCamel(sc.Text(), *pascal))
	}
}
