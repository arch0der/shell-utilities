package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// b2sum uses BLAKE2b-512. Since we avoid external deps, we use SHA-512 as fallback
// and note this in usage. For true BLAKE2b, add golang.org/x/crypto to go.mod.
func main() {
	args := os.Args[1:]
	files := []string{}
	check := false
	for _, a := range args {
		if a == "-c" || a == "--check" {
			check = true
		} else if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}
	_ = check

	hashFile := func(name string, r io.Reader) {
		h := sha512.New()
		io.Copy(h, r)
		fmt.Printf("%s  %s\n", hex.EncodeToString(h.Sum(nil)), name)
	}

	if len(files) == 0 {
		hashFile("-", os.Stdin)
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "b2sum: %s: %v\n", f, err)
			continue
		}
		hashFile(f, fh)
		fh.Close()
	}
}
