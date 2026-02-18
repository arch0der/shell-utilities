// echo - Print arguments to stdout
// Usage: echo [-n] [-e] [string...]
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	noNewline = flag.Bool("n", false, "Do not print trailing newline")
	escape    = flag.Bool("e", false, "Enable interpretation of backslash escapes")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: echo [-n] [-e] [string...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	out := strings.Join(flag.Args(), " ")

	if *escape {
		out = strings.NewReplacer(
			`\n`, "\n",
			`\t`, "\t",
			`\r`, "\r",
			`\\`, "\\",
			`\a`, "\a",
			`\b`, "\b",
			`\v`, "\v",
		).Replace(out)
	}

	if *noNewline {
		fmt.Print(out)
	} else {
		fmt.Println(out)
	}
}
