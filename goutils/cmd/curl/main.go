// curl - Make HTTP requests
// Usage: curl [-X method] [-H header] [-d data] [-o file] [-i] <url>
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type headers []string

func (h *headers) String() string  { return strings.Join(*h, ", ") }
func (h *headers) Set(v string) error { *h = append(*h, v); return nil }

var (
	method  = flag.String("X", "GET", "HTTP method")
	data    = flag.String("d", "", "Request body data")
	output  = flag.String("o", "", "Output file (default: stdout)")
	include = flag.Bool("i", false, "Include response headers in output")
	hdrs    headers
)

func main() {
	flag.Var(&hdrs, "H", "HTTP header (can be used multiple times)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: curl [-X method] [-H header] [-d data] [-o file] [-i] <url>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	url := flag.Arg(0)

	var body io.Reader
	if *data != "" {
		body = strings.NewReader(*data)
		if *method == "GET" {
			*method = "POST"
		}
	}

	req, err := http.NewRequest(*method, url, body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "curl:", err)
		os.Exit(1)
	}

	for _, h := range hdrs {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "curl:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var out io.Writer = os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintln(os.Stderr, "curl:", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	if *include {
		fmt.Fprintf(out, "HTTP/1.1 %s\n", resp.Status)
		for k, vs := range resp.Header {
			for _, v := range vs {
				fmt.Fprintf(out, "%s: %s\n", k, v)
			}
		}
		fmt.Fprintln(out)
	}

	io.Copy(out, resp.Body)
}
