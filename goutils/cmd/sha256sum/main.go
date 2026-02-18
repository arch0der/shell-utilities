// sha256sum - Compute SHA-256 checksums
// Usage: sha256sum [-c] [file...]
package main

import (
	"bufio"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var check = flag.Bool("c", false, "Check SHA-256 sums from file")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: sha256sum [-c] [file...]")
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
					fmt.Fprintln(os.Stderr, "sha256sum:", err)
					continue
				}
				checkSums(f)
				f.Close()
			}
		}
		return
	}

	if flag.NArg() == 0 {
		sum, _ := hashReader(os.Stdin)
		fmt.Printf("%s  -\n", sum)
		return
	}

	exitCode := 0
	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "sha256sum:", err)
			exitCode = 1
			continue
		}
		sum, err := hashReader(f)
		f.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "sha256sum:", err)
			exitCode = 1
			continue
		}
		fmt.Printf("%s  %s\n", sum, path)
	}
	os.Exit(exitCode)
}

func hashReader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func checkSums(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "  ", 2)
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
