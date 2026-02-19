package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	decode := false
	wrap := 76
	files := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-d" || a == "--decode":
			decode = true
		case a == "-w" && i+1 < len(args):
			i++
			fmt.Sscan(args[i], &wrap)
		case a == "-w0":
			wrap = 0
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}
	var r io.Reader = os.Stdin
	if len(files) > 0 {
		f, err := os.Open(files[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "base64: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}
	data, _ := io.ReadAll(r)
	if decode {
		clean := strings.ReplaceAll(strings.ReplaceAll(string(data), "\n", ""), "\r", "")
		out, err := base64.StdEncoding.DecodeString(clean)
		if err != nil {
			fmt.Fprintf(os.Stderr, "base64: invalid input\n")
			os.Exit(1)
		}
		os.Stdout.Write(out)
	} else {
		enc := base64.StdEncoding.EncodeToString(data)
		w := bufio.NewWriter(os.Stdout)
		if wrap == 0 {
			w.WriteString(enc + "\n")
		} else {
			for len(enc) > wrap {
				w.WriteString(enc[:wrap] + "\n")
				enc = enc[wrap:]
			}
			if len(enc) > 0 {
				w.WriteString(enc + "\n")
			}
		}
		w.Flush()
	}
}
