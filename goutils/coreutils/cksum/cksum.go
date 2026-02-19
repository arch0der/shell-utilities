package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

// CRC-32 used by POSIX cksum
var cksumTable [256]uint32

func init() {
	for i := range cksumTable {
		crc := uint32(i) << 24
		for j := 0; j < 8; j++ {
			if crc&0x80000000 != 0 {
				crc = (crc << 1) ^ 0x04C11DB7
			} else {
				crc <<= 1
			}
		}
		cksumTable[i] = crc
	}
}

func cksumCompute(data []byte) uint32 {
	crc := uint32(0)
	for _, b := range data {
		crc = (crc << 8) ^ cksumTable[byte(crc>>24)^b]
	}
	// include length bytes
	length := uint64(len(data))
	var lb [8]byte
	binary.BigEndian.PutUint64(lb[:], length)
	start := 0
	for start < 7 && lb[start] == 0 {
		start++
	}
	for _, b := range lb[start:] {
		crc = (crc << 8) ^ cksumTable[byte(crc>>24)^b]
	}
	return ^crc
}

func main() {
	args := os.Args[1:]
	files := []string{}
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}
	process := func(name string, r io.Reader) {
		data, _ := io.ReadAll(r)
		crc := cksumCompute(data)
		fmt.Printf("%d %d %s\n", crc, len(data), name)
	}
	if len(files) == 0 {
		process("", os.Stdin)
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cksum: %s: %v\n", f, err)
			continue
		}
		process(f, fh)
		fh.Close()
	}
}
