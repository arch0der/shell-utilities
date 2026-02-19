// urlencode - URL-encode strings
// Usage: urlencode [-d] [string...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

var decode = flag.Bool("d", false, "Decode instead of encode")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: urlencode [-d] [string...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	process := func(s string) {
		if *decode {
			decoded, err := url.QueryUnescape(s)
			if err != nil {
				fmt.Fprintln(os.Stderr, "urlencode:", err)
				return
			}
			fmt.Println(decoded)
		} else {
			fmt.Println(url.QueryEscape(s))
		}
	}

	if flag.NArg() > 0 {
		process(strings.Join(flag.Args(), " "))
		return
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		process(sc.Text())
	}
	_ = io.EOF // avoid unused import
}
