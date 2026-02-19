// checksum - compute MD5, SHA1, SHA256, SHA512, CRC32 simultaneously
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

func sum(r io.Reader, name string) {
	h1, h2, h3, h4, h5 := md5.New(), sha1.New(), sha256.New(), sha512.New(), crc32.NewIEEE()
	if _, err := io.Copy(io.MultiWriter(h1, h2, h3, h4, h5), r); err != nil {
		fmt.Fprintln(os.Stderr, err); return
	}
	fmt.Printf("File    : %s\n", name)
	fmt.Printf("MD5     : %s\n", hex.EncodeToString(h1.Sum(nil)))
	fmt.Printf("SHA1    : %s\n", hex.EncodeToString(h2.Sum(nil)))
	fmt.Printf("SHA256  : %s\n", hex.EncodeToString(h3.Sum(nil)))
	fmt.Printf("SHA512  : %s\n", hex.EncodeToString(h4.Sum(nil)))
	fmt.Printf("CRC32   : %08X\n\n", h5.Sum32())
}

func main() {
	if len(os.Args) == 1 { sum(os.Stdin, "<stdin>"); return }
	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		sum(f, path); f.Close()
	}
}
