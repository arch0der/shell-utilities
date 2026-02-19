package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func init() { register("split", runSplit) }

func runSplit() {
	args := os.Args[1:]
	linesPerFile := 1000
	bytesPerFile := int64(-1)
	prefix := "x"
	suffix := "aa"
	numericSuffix := false
	verbose := false
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-l" && i+1 < len(args):
			i++
			linesPerFile, _ = strconv.Atoi(args[i])
		case strings.HasPrefix(a, "-l"):
			linesPerFile, _ = strconv.Atoi(a[2:])
		case a == "-b" && i+1 < len(args):
			i++
			bytesPerFile = parseSize(args[i])
		case strings.HasPrefix(a, "-b"):
			bytesPerFile = parseSize(a[2:])
		case a == "-d" || a == "--numeric-suffixes":
			numericSuffix = true
		case a == "-v" || a == "--verbose":
			verbose = true
		case a == "--suffix-length" && i+1 < len(args):
			i++
			// ignored for simplicity
		case !strings.HasPrefix(a, "-"):
			if len(files) == 0 {
				files = append(files, a)
			} else {
				prefix = a
			}
		}
	}
	_ = suffix

	var r io.Reader = os.Stdin
	if len(files) > 0 {
		fh, err := os.Open(files[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "split: %v\n", err)
			os.Exit(1)
		}
		defer fh.Close()
		r = fh
	}

	getFilename := func(n int) string {
		if numericSuffix {
			return fmt.Sprintf("%s%02d", prefix, n)
		}
		// aa, ab, ac, ...
		letters := "abcdefghijklmnopqrstuvwxyz"
		suf := string([]byte{letters[n/26 % 26], letters[n % 26]})
		return prefix + suf
	}

	fileNum := 0
	var curFile *os.File
	var curSize int64
	var curLines int

	newFile := func() {
		if curFile != nil {
			curFile.Close()
		}
		name := getFilename(fileNum)
		if verbose {
			fmt.Fprintf(os.Stderr, "creating file '%s'\n", name)
		}
		var err error
		curFile, err = os.Create(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "split:", err)
			os.Exit(1)
		}
		fileNum++
		curSize = 0
		curLines = 0
	}

	newFile()
	if bytesPerFile > 0 {
		buf := make([]byte, 4096)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				remaining := n
				for remaining > 0 {
					canWrite := int(bytesPerFile - curSize)
					if canWrite <= 0 {
						newFile()
						canWrite = int(bytesPerFile)
					}
					toWrite := remaining
					if toWrite > canWrite {
						toWrite = canWrite
					}
					curFile.Write(buf[:toWrite])
					curSize += int64(toWrite)
					remaining -= toWrite
				}
			}
			if err != nil {
				break
			}
		}
	} else {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for sc.Scan() {
			if curLines >= linesPerFile {
				newFile()
			}
			curFile.WriteString(sc.Text() + "\n")
			curLines++
		}
	}
	if curFile != nil {
		curFile.Close()
	}
}
