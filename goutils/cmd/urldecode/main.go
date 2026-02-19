// urldecode - URL decode strings
// Usage: urldecode [string...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var query = flag.Bool("q", false, "Decode query string (+ as space)")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: urldecode [-q] [string...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	process := func(s string) {
		var decoded string
		var err error
		if *query {
			decoded, err = url.QueryUnescape(s)
		} else {
			decoded, err = url.PathUnescape(s)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "urldecode:", err)
			return
		}
		fmt.Println(decoded)
	}

	if flag.NArg() > 0 {
		process(strings.Join(flag.Args(), " "))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		process(sc.Text())
	}
}
