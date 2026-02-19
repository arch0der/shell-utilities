package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func init() { register("md5sum", runMd5sum) }

func hashWithAlgo(name string, r io.Reader) string {
	h := md5.New()
	switch name {
	case "md5sum":
		h = md5.New()
	}
	io.Copy(h, r)
	return hex.EncodeToString(h.Sum(nil))
}

func runHashCmd(algo string) {
	args := os.Args[1:]
	check := false
	tag := false
	files := []string{}
	for _, a := range args {
		switch a {
		case "-c", "--check":
			check = true
		case "--tag":
			tag = true
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}
	_ = check
	_ = tag

	doHash := func(r io.Reader) string {
		return hashWithAlgo(algo, r)
	}

	if len(files) == 0 {
		sum := doHash(os.Stdin)
		fmt.Printf("%s  -\n", sum)
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %v\n", algo, f, err)
			continue
		}
		sum := doHash(fh)
		fh.Close()
		fmt.Printf("%s  %s\n", sum, f)
	}
}

func runMd5sum() {
	args := os.Args[1:]
	check := false
	files := []string{}
	for _, a := range args {
		if a == "-c" || a == "--check" {
			check = true
		} else if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}
	if check && len(files) > 0 {
		exitCode := 0
		for _, f := range files {
			fh, _ := os.Open(f)
			sc := bufio.NewScanner(fh)
			for sc.Scan() {
				line := sc.Text()
				parts := strings.SplitN(line, "  ", 2)
				if len(parts) != 2 {
					continue
				}
				expected, fname := parts[0], parts[1]
				data, err := os.ReadFile(fname)
				if err != nil {
					fmt.Fprintf(os.Stderr, "md5sum: %s: %v\n", fname, err)
					exitCode = 1
					continue
				}
				h := md5.Sum(data)
				actual := hex.EncodeToString(h[:])
				if actual == expected {
					fmt.Printf("%s: OK\n", fname)
				} else {
					fmt.Printf("%s: FAILED\n", fname)
					exitCode = 1
				}
			}
			fh.Close()
		}
		os.Exit(exitCode)
	}
	if len(files) == 0 {
		h := md5.New()
		io.Copy(h, os.Stdin)
		fmt.Printf("%s  -\n", hex.EncodeToString(h.Sum(nil)))
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "md5sum: %s: %v\n", f, err)
			continue
		}
		h := md5.New()
		io.Copy(h, fh)
		fh.Close()
		fmt.Printf("%s  %s\n", hex.EncodeToString(h.Sum(nil)), f)
	}
}
