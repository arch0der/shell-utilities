// grep - Search text using patterns
// Usage: grep [-i] [-n] [-v] [-r] <pattern> [file...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ignoreCase  = flag.Bool("i", false, "Case-insensitive matching")
	lineNumbers = flag.Bool("n", false, "Show line numbers")
	invert      = flag.Bool("v", false, "Invert match (select non-matching lines)")
	recursive   = flag.Bool("r", false, "Recursively search directories")
	count       = flag.Bool("c", false, "Print only count of matching lines")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: grep [-i] [-n] [-v] [-r] [-c] <pattern> [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	pattern := flag.Arg(0)
	if *ignoreCase {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "grep: invalid pattern:", err)
		os.Exit(1)
	}

	files := flag.Args()[1:]
	if len(files) == 0 {
		searchReader(os.Stdin, "", re)
		return
	}

	multiFile := len(files) > 1
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "grep:", err)
			continue
		}
		if info.IsDir() && *recursive {
			filepath.Walk(path, func(p string, i os.FileInfo, e error) error {
				if e != nil || i.IsDir() {
					return nil
				}
				f, err := os.Open(p)
				if err != nil {
					return nil
				}
				defer f.Close()
				searchReader(f, p, re)
				return nil
			})
		} else {
			f, err := os.Open(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "grep:", err)
				continue
			}
			prefix := ""
			if multiFile {
				prefix = path
			}
			searchReader(f, prefix, re)
			f.Close()
		}
	}
}

func searchReader(r *os.File, prefix string, re *regexp.Regexp) {
	scanner := bufio.NewScanner(r)
	lineNo := 0
	matched := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		match := re.MatchString(line)
		if *invert {
			match = !match
		}
		if match {
			matched++
			if !*count {
				out := ""
				if prefix != "" {
					out += prefix + ":"
				}
				if *lineNumbers {
					out += fmt.Sprintf("%d:", lineNo)
				}
				out += line
				fmt.Println(out)
			}
		}
	}
	if *count {
		if prefix != "" {
			fmt.Printf("%s:%d\n", prefix, matched)
		} else {
			fmt.Println(matched)
		}
	}
}
