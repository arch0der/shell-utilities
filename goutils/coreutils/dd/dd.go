package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	inFile := ""
	outFile := ""
	bs := int64(512)
	ibs := int64(512)
	obs := int64(512)
	count := int64(-1)
	skip := int64(0)
	seek := int64(0)
	conv := ""

	for _, a := range args {
		kv := strings.SplitN(a, "=", 2)
		if len(kv) != 2 {
			continue
		}
		k, v := kv[0], kv[1]
		switch k {
		case "if":
			inFile = v
		case "of":
			outFile = v
		case "bs":
			bs = parseSize(v)
			ibs = bs
			obs = bs
		case "ibs":
			ibs = parseSize(v)
		case "obs":
			obs = parseSize(v)
		case "count":
			count, _ = strconv.ParseInt(v, 10, 64)
		case "skip":
			skip = parseSize(v)
		case "seek":
			seek = parseSize(v)
		case "conv":
			conv = v
		}
	}
	_ = bs
	_ = obs
	_ = conv

	var in io.Reader = os.Stdin
	var out io.Writer = os.Stdout

	if inFile != "" {
		f, err := os.Open(inFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dd: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		if skip > 0 {
			f.Seek(skip*ibs, io.SeekStart)
		}
		in = f
	}

	if outFile != "" {
		f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dd: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		if seek > 0 {
			f.Seek(seek*obs, io.SeekStart)
		}
		out = f
	}

	buf := make([]byte, ibs)
	var blocks, bytes int64
	for count < 0 || blocks < count {
		n, err := in.Read(buf)
		if n > 0 {
			if strings.Contains(conv, "ucase") {
				for i := range buf[:n] {
					if buf[i] >= 'a' && buf[i] <= 'z' {
						buf[i] -= 32
					}
				}
			} else if strings.Contains(conv, "lcase") {
				for i := range buf[:n] {
					if buf[i] >= 'A' && buf[i] <= 'Z' {
						buf[i] += 32
					}
				}
			}
			out.Write(buf[:n])
			bytes += int64(n)
			blocks++
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "dd: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Fprintf(os.Stderr, "%d+0 records in\n%d+0 records out\n%d bytes transferred\n", blocks, blocks, bytes)
}
