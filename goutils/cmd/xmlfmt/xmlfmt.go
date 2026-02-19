// xmlfmt - pretty-print or minify XML
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func prettyPrint(data []byte, indent string) (string, error) {
	var buf bytes.Buffer
	decoder := xml.NewDecoder(bytes.NewReader(data))
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", indent)
	for {
		token, err := decoder.Token()
		if err == io.EOF { break }
		if err != nil { return "", err }
		if err := encoder.EncodeToken(token); err != nil { return "", err }
	}
	if err := encoder.Flush(); err != nil { return "", err }
	return buf.String(), nil
}

func minify(data []byte) (string, error) {
	var buf bytes.Buffer
	decoder := xml.NewDecoder(bytes.NewReader(data))
	encoder := xml.NewEncoder(&buf)
	for {
		token, err := decoder.Token()
		if err == io.EOF { break }
		if err != nil { return "", err }
		if cd, ok := token.(xml.CharData); ok {
			trimmed := strings.TrimSpace(string(cd))
			if trimmed == "" { continue }
			token = xml.CharData(trimmed)
		}
		if err := encoder.EncodeToken(token); err != nil { return "", err }
	}
	encoder.Flush()
	return buf.String(), nil
}

func main() {
	doMinify := false
	indent := "  "
	files := []string{}
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-m", "--minify": doMinify = true
		case "-i": i++; indent = os.Args[i]
		default: files = append(files, os.Args[i])
		}
	}

	process := func(r io.Reader, name string) {
		data, err := io.ReadAll(r)
		if err != nil { fmt.Fprintln(os.Stderr, err); return }
		var out string
		if doMinify { out, err = minify(data) } else { out, err = prettyPrint(data, indent) }
		if err != nil { fmt.Fprintf(os.Stderr, "xmlfmt: %s: %v\n", name, err); return }
		fmt.Println(out)
	}

	if len(files) == 0 { process(os.Stdin, "<stdin>"); return }
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		process(fh, f); fh.Close()
	}
}
