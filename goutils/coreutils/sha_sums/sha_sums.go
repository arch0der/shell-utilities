package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

func init() {

}

func runHashTool(name string, newHash func() hash.Hash) {
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
					fmt.Fprintf(os.Stderr, "%s: %s: %v\n", name, fname, err)
					exitCode = 1
					continue
				}
				h := newHash()
				h.Write(data)
				actual := hex.EncodeToString(h.Sum(nil))
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
		h := newHash()
		io.Copy(h, os.Stdin)
		fmt.Printf("%s  -\n", hex.EncodeToString(h.Sum(nil)))
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %v\n", name, f, err)
			continue
		}
		h := newHash()
		io.Copy(h, fh)
		fh.Close()
		fmt.Printf("%s  %s\n", hex.EncodeToString(h.Sum(nil)), f)
	}
}

func main() { runHashTool("sha1sum", sha1.New) }
func main() { runHashTool("sha224sum", sha256.New224) }
func main() { runHashTool("sha256sum", sha256.New) }
func main() { runHashTool("sha384sum", sha512.New384) }
func main() { runHashTool("sha512sum", sha512.New) }
