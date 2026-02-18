// md5sum - Compute MD5 checksums
// Usage: md5sum [-c] [file...]
package main

import (
	"bufio"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var check = flag.Bool("c", false, "Check MD5 sums from file")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: md5sum [-c] [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *check {
		if flag.NArg() == 0 {
			checkSums(os.Stdin)
		} else {
			for _, path := range flag.Args() {
				f, err := os.Open(path)
				if err != nil {
					fmt.Fprintln(os.Stderr, "md5sum:", err)
					continue
				}
				checkSums(f)
				f.Close()
			}
		}
		return
	}

	if flag.NArg() == 0 {
		sum, err := hashReader(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "md5sum:", err)
			os.Exit(1)
		}
		fmt.Printf("%s  -\n", sum)
		return
	}

	exitCode := 0
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "md5sum:", err)
			exitCode = 1
			continue
		}
		sum, err := hashReader(f)
		f.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "md5sum:", err)
			exitCode = 1
			continue
		}
		fmt.Printf("%s  %s\n", sum, path)
	}
	os.Exit(exitCode)
}

func hashReader(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func checkSums(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) != 2 {
			continue
		}
		expected, path := parts[0], parts[1]
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("%s: FAILED open or read\n", path)
			continue
		}
		got, err := hashReader(f)
		f.Close()
		if err != nil || got != expected {
			fmt.Printf("%s: FAILED\n", path)
		} else {
			fmt.Printf("%s: OK\n", path)
		}
	}
}
