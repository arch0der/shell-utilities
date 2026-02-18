// mkdir - Create directories
// Usage: mkdir [-p] <dir>...
package main

import (
	"flag"
	"fmt"
	"os"
)

var parents = flag.Bool("p", false, "Create parent directories as needed")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: mkdir [-p] <dir>...")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	exitCode := 0
	for _, dir := range flag.Args() {
		var err error
		if *parents {
			err = os.MkdirAll(dir, 0755)
		} else {
			err = os.Mkdir(dir, 0755)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "mkdir:", err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
