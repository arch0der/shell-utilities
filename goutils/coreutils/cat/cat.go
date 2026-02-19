package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	showNum := false
	showEnds := false
	showTabs := false
	showNonprint := false
	squeezeBlank := false
	files := []string{}

	for _, a := range args {
		switch a {
		case "-n", "--number":
			showNum = true
		case "-b", "--number-nonblank":
			showNum = true
		case "-A", "--show-all":
			showEnds, showTabs, showNonprint = true, true, true
		case "-e":
			showEnds, showNonprint = true, true
		case "-t":
			showTabs, showNonprint = true, true
		case "-v", "--show-nonprinting":
			showNonprint = true
		case "-E", "--show-ends":
			showEnds = true
		case "-T", "--show-tabs":
			showTabs = true
		case "-s", "--squeeze-blank":
			squeezeBlank = true
		default:
			if !strings.HasPrefix(a, "-") || a == "-" {
				files = append(files, a)
			}
		}
	}

	lineNum := 1
	lastBlank := false
	catReader := func(r io.Reader) {
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 1<<20), 1<<20)
		for scanner.Scan() {
			line := scanner.Text()
			if squeezeBlank && line == "" {
				if lastBlank {
					continue
				}
				lastBlank = true
			} else {
				lastBlank = false
			}
			if showNonprint {
				var sb strings.Builder
				for i := 0; i < len(line); i++ {
					c := line[i]
					if showTabs && c == '\t' {
						sb.WriteString("^I")
					} else if c >= 32 && c < 127 {
						sb.WriteByte(c)
					} else if c == 127 {
						sb.WriteString("^?")
					} else if c < 32 {
						sb.WriteByte('^')
						sb.WriteByte(c + 64)
					} else {
						sb.WriteString("M-")
						if c-128 < 32 {
							sb.WriteByte('^')
							sb.WriteByte(c - 128 + 64)
						} else {
							sb.WriteByte(c - 128)
						}
					}
				}
				line = sb.String()
			} else if showTabs {
				line = strings.ReplaceAll(line, "\t", "^I")
			}
			if showNum {
				fmt.Printf("%6d\t", lineNum)
				lineNum++
			}
			if showEnds {
				fmt.Println(line + "$")
			} else {
				fmt.Println(line)
			}
		}
	}

	if len(files) == 0 {
		catReader(os.Stdin)
		return
	}
	for _, f := range files {
		if f == "-" {
			catReader(os.Stdin)
			continue
		}
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cat: %s: %v\n", f, err)
			continue
		}
		catReader(fh)
		fh.Close()
	}
}
