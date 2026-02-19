// htmlstrip - remove HTML tags from input, optionally decoding entities
package main

import (
	"bufio"
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"
)

var (
	tagRe    = regexp.MustCompile(`<[^>]+>`)
	scriptRe = regexp.MustCompile(`(?is)<(script|style)[^>]*>.*?</(script|style)>`)
	spaceRe  = regexp.MustCompile(`[ \t]+`)
)

func strip(s string, decode, collapseSpace bool) string {
	s = scriptRe.ReplaceAllString(s, "")
	s = tagRe.ReplaceAllString(s, "")
	if decode { s = html.UnescapeString(s) }
	if collapseSpace {
		s = spaceRe.ReplaceAllString(s, " ")
		s = strings.TrimSpace(s)
	}
	return s
}

func main() {
	decode := false
	collapse := false
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-d", "--decode": decode = true
		case "-c", "--collapse": collapse = true
		}
	}
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() { fmt.Println(strip(sc.Text(), decode, collapse)) }
}
