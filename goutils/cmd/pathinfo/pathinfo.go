// pathinfo - dissect file paths into components
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: pathinfo <path> [part]")
		fmt.Fprintln(os.Stderr, "  parts: dir | base | ext | stem | abs | clean")
		os.Exit(1)
	}
	path := os.Args[1]
	part := ""
	if len(os.Args) > 2 { part = os.Args[2] }

	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(path)
	stem := strings.TrimSuffix(base, ext)
	abs, _ := filepath.Abs(path)
	clean := filepath.Clean(path)

	switch part {
	case "dir": fmt.Println(dir)
	case "base": fmt.Println(base)
	case "ext": fmt.Println(ext)
	case "stem": fmt.Println(stem)
	case "abs": fmt.Println(abs)
	case "clean": fmt.Println(clean)
	case "": // print all
		fmt.Printf("Path    : %s\n", path)
		fmt.Printf("Dir     : %s\n", dir)
		fmt.Printf("Base    : %s\n", base)
		fmt.Printf("Stem    : %s\n", stem)
		fmt.Printf("Ext     : %s\n", ext)
		fmt.Printf("Abs     : %s\n", abs)
		fmt.Printf("Clean   : %s\n", clean)
	default:
		fmt.Fprintf(os.Stderr, "pathinfo: unknown part %q\n", part); os.Exit(1)
	}
}
