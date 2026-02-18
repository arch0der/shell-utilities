// ln - Create hard or symbolic links
// Usage: ln [-s] [-f] <target> <link>
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	symbolic = flag.Bool("s", false, "Create symbolic link")
	force    = flag.Bool("f", false, "Remove existing destination")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: ln [-s] [-f] <target> <link>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	target := flag.Arg(0)
	link := flag.Arg(1)

	// If link is an existing directory, put link inside it
	if info, err := os.Stat(link); err == nil && info.IsDir() {
		link = filepath.Join(link, filepath.Base(target))
	}

	if *force {
		os.Remove(link)
	}

	var err error
	if *symbolic {
		err = os.Symlink(target, link)
	} else {
		err = os.Link(target, link)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "ln:", err)
		os.Exit(1)
	}
}
