// snake_case - Convert strings to snake_case, kebab-case, or SCREAMING_SNAKE.
//
// Usage:
//
//	snake_case [OPTIONS] [TEXT...]
//	echo "HelloWorld" | snake_case
//
// Options:
//
//	-k        kebab-case (use - instead of _)
//	-u        SCREAMING_SNAKE_CASE (uppercase)
//	-s SEP    Output separator (overrides -k)
//
// Examples:
//
//	echo "HelloWorld"     | snake_case    # hello_world
//	echo "myCSVParser"    | snake_case    # my_csv_parser
//	echo "Hello World"    | snake_case -k # hello-world
//	echo "hello_world"    | snake_case -u # HELLO_WORLD
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
	kebab  = flag.Bool("k", false, "kebab-case")
	scream = flag.Bool("u", false, "SCREAMING_SNAKE")
	sep    = flag.String("s", "", "output separator")
)

func toSnake(s string) string {
	separator := "_"
	if *sep != "" {
		separator = *sep
	} else if *kebab {
		separator = "-"
	}

	var words []string
	var cur strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		if r == '_' || r == '-' || r == '.' || r == ' ' {
			if cur.Len() > 0 {
				words = append(words, cur.String())
				cur.Reset()
			}
		} else if unicode.IsUpper(r) {
			if cur.Len() > 0 {
				// Check for acronym: AABcc → AA + Bcc
				if i+1 < len(runes) && unicode.IsUpper(runes[i+1]) && cur.Len() > 1 && unicode.IsLower(runes[i-1]) {
					words = append(words, cur.String())
					cur.Reset()
				} else if unicode.IsLower(runes[i-1]) {
					words = append(words, cur.String())
					cur.Reset()
				}
			}
			cur.WriteRune(unicode.ToLower(r))
		} else {
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		words = append(words, cur.String())
	}

	result := strings.Join(words, separator)
	if *scream {
		return strings.ToUpper(result)
	}
	return result
}

func main() {
	flag.Parse()
	args := flag.Args()
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if len(args) > 0 {
		fmt.Fprintln(w, toSnake(strings.Join(args, " ")))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Fprintln(w, toSnake(sc.Text()))
	}
}
