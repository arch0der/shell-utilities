// mktemp - Create temporary files or directories
// Usage: mktemp [-d] [-p dir] [template]
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	dir      = flag.Bool("d", false, "Create a temporary directory")
	tempDir  = flag.String("p", "", "Use specified directory as base (default: $TMPDIR)")
	suffix   = flag.String("suffix", "", "Append suffix to template")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: mktemp [-d] [-p dir] [--suffix suf] [template]")
		flag.PrintDefaults()
	}
	flag.Parse()

	prefix := "tmp"
	if flag.NArg() > 0 {
		prefix = flag.Arg(0)
	}

	base := *tempDir
	if base == "" {
		base = os.TempDir()
	}

	pattern := prefix
	if *suffix != "" {
		pattern += "*" + *suffix
	}

	if *dir {
		d, err := os.MkdirTemp(base, pattern)
		if err != nil {
			fmt.Fprintln(os.Stderr, "mktemp:", err)
			os.Exit(1)
		}
		fmt.Println(d)
	} else {
		f, err := os.CreateTemp(base, pattern)
		if err != nil {
			fmt.Fprintln(os.Stderr, "mktemp:", err)
			os.Exit(1)
		}
		f.Close()
		fmt.Println(f.Name())
	}
}
