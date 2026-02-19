package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	directory := false
	dryRun := false
	quiet := false
	tmpdir := os.TempDir()
	suffix := ""
	template := ""

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-d" || a == "--directory":
			directory = true
		case a == "-u" || a == "--dry-run":
			dryRun = true
		case a == "-q" || a == "--quiet":
			quiet = true
		case a == "-p" && i+1 < len(args):
			i++
			tmpdir = args[i]
		case strings.HasPrefix(a, "--tmpdir="):
			tmpdir = a[9:]
		case a == "--tmpdir":
			if i+1 < len(args) {
				i++
				tmpdir = args[i]
			}
		case a == "--suffix" && i+1 < len(args):
			i++
			suffix = args[i]
		case strings.HasPrefix(a, "--suffix="):
			suffix = a[9:]
		case !strings.HasPrefix(a, "-"):
			template = a
		}
	}

	if template == "" {
		template = "tmp.XXXXXXXXXX"
	}
	if suffix != "" && !strings.HasSuffix(template, "XXXXXX") {
		template = template + suffix
	}

	// Replace template XXXXXX with temp file
	pattern := strings.TrimSuffix(template, strings.Repeat("X", strings.Count(template, "X")))
	_ = pattern

	if directory {
		dir, err := os.MkdirTemp(tmpdir, template)
		if err != nil {
			if !quiet {
				fmt.Fprintln(os.Stderr, "mktemp:", err)
			}
			os.Exit(1)
		}
		if dryRun {
			os.Remove(dir)
		}
		fmt.Println(dir)
	} else {
		f, err := os.CreateTemp(tmpdir, template)
		if err != nil {
			if !quiet {
				fmt.Fprintln(os.Stderr, "mktemp:", err)
			}
			os.Exit(1)
		}
		name := f.Name()
		f.Close()
		if dryRun {
			os.Remove(name)
		}
		fmt.Println(name)
	}
}
