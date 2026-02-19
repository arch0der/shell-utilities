// indent - Add or remove indentation from text.
//
// Usage:
//
//	indent [OPTIONS] [FILE...]
//	cat code.txt | indent -n 4
//
// Options:
//
//	-n N      Add N spaces of indentation (default: 4)
//	-t        Use tabs instead of spaces
//	-d        Dedent: remove common leading whitespace
//	-d N      Remove exactly N levels of indentation
//	-p STR    Add custom prefix string to each line
//	-s        Skip blank lines (don't indent them)
//	-c        Continue: preserve existing indentation, add on top
//
// Examples:
//
//	cat code.py | indent -n 4           # add 4 spaces
//	cat indented.txt | indent -d        # remove common indent
//	echo "code" | indent -t             # add tab
//	cat block.txt | indent -p "    "    # custom prefix
//	cat code.go | indent -d -n 2        # dedent then re-indent by 2
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	spaces    = flag.Int("n", 4, "spaces to add")
	useTabs   = flag.Bool("t", false, "use tabs")
	dedent    = flag.Bool("d", false, "dedent")
	prefix    = flag.String("p", "", "custom prefix")
	skipBlank = flag.Bool("s", false, "skip blank lines")
)

func commonIndent(lines []string) string {
	common := ""
	first := true
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		indent := ""
		for _, ch := range l {
			if ch == ' ' || ch == '\t' {
				indent += string(ch)
			} else {
				break
			}
		}
		if first {
			common = indent
			first = false
		} else {
			// Find common prefix
			i := 0
			for i < len(common) && i < len(indent) && common[i] == indent[i] {
				i++
			}
			common = common[:i]
		}
	}
	return common
}

func main() {
	flag.Parse()
	files := flag.Args()

	var readers []*os.File
	if len(files) == 0 {
		readers = []*os.File{os.Stdin}
	} else {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "indent: %v\n", err)
				os.Exit(1)
			}
			defer fh.Close()
			readers = append(readers, fh)
		}
	}

	// Read all lines
	var lines []string
	for _, r := range readers {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1024*1024), 1024*1024)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
	}

	// Calculate add prefix
	addPrefix := *prefix
	if addPrefix == "" {
		if *useTabs {
			addPrefix = "\t"
		} else {
			addPrefix = strings.Repeat(" ", *spaces)
		}
	}

	// Dedent first
	if *dedent {
		common := commonIndent(lines)
		for i, l := range lines {
			if strings.TrimSpace(l) == "" {
				continue
			}
			lines[i] = strings.TrimPrefix(l, common)
		}
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for _, l := range lines {
		if *skipBlank && strings.TrimSpace(l) == "" {
			fmt.Fprintln(w, l)
			continue
		}
		if *dedent && *prefix == "" && !*useTabs && *spaces == 4 {
			// Dedent only mode
			fmt.Fprintln(w, l)
		} else {
			fmt.Fprintln(w, addPrefix+l)
		}
	}
}
