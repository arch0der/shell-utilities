// slugify - Convert text to URL-safe slugs.
//
// Usage:
//
//	slugify [OPTIONS] [TEXT...]
//	echo "Hello World" | slugify
//
// Options:
//
//	-s SEP    Separator (default: -)
//	-u        Uppercase output
//	--keep    Extra chars to preserve (default: "")
//
// Examples:
//
//	slugify "Hello, World!"           # hello-world
//	slugify -s _ "My File Name"       # my_file_name
//	echo "Héllo Wörld" | slugify      # hello-world
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var sep = flag.String("s", "-", "separator")
var upper = flag.Bool("u", false, "uppercase")

func toASCII(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func slugify(s string) string {
	s = toASCII(s)
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, *sep)
	s = strings.Trim(s, *sep)
	if *upper {
		s = strings.ToUpper(s)
	}
	return s
}

func main() {
	flag.Parse()
	args := flag.Args()
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if len(args) > 0 {
		fmt.Fprintln(w, slugify(strings.Join(args, " ")))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Fprintln(w, slugify(sc.Text()))
	}
}
