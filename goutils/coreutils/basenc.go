package main

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func init() { register("basenc", runBasenc) }

func runBasenc() {
	args := os.Args[1:]
	decode := false
	encoding := "base64"
	files := []string{}
	for _, a := range args {
		switch a {
		case "-d", "--decode":
			decode = true
		case "--base64":
			encoding = "base64"
		case "--base64url":
			encoding = "base64url"
		case "--base32":
			encoding = "base32"
		case "--base32hex":
			encoding = "base32hex"
		case "--base16":
			encoding = "base16"
		case "--base2msbf":
			encoding = "base2msbf"
		case "--z85":
			encoding = "z85"
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}
	var r io.Reader = os.Stdin
	if len(files) > 0 {
		f, _ := os.Open(files[0])
		defer f.Close()
		r = f
	}
	data, _ := io.ReadAll(r)
	if decode {
		var out []byte
		var err error
		clean := strings.ReplaceAll(string(data), "\n", "")
		switch encoding {
		case "base64":
			out, err = base64.StdEncoding.DecodeString(clean)
		case "base64url":
			out, err = base64.URLEncoding.DecodeString(clean)
		case "base32":
			out, err = base32.StdEncoding.DecodeString(clean)
		case "base32hex":
			out, err = base32.HexEncoding.DecodeString(clean)
		case "base16":
			out, err = hex.DecodeString(clean)
		default:
			fmt.Fprintf(os.Stderr, "basenc: unsupported encoding %s\n", encoding)
			os.Exit(1)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "basenc: decode error: %v\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(out)
	} else {
		var enc string
		switch encoding {
		case "base64":
			enc = base64.StdEncoding.EncodeToString(data)
		case "base64url":
			enc = base64.URLEncoding.EncodeToString(data)
		case "base32":
			enc = base32.StdEncoding.EncodeToString(data)
		case "base32hex":
			enc = base32.HexEncoding.EncodeToString(data)
		case "base16":
			enc = strings.ToUpper(hex.EncodeToString(data))
		default:
			fmt.Fprintf(os.Stderr, "basenc: unsupported encoding %s\n", encoding)
			os.Exit(1)
		}
		fmt.Println(enc)
	}
}
