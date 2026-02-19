// headers - Show HTTP response headers for a URL.
//
// Usage:
//
//	headers [OPTIONS] URL
//
// Options:
//
//	-m METHOD   HTTP method (default: HEAD)
//	-H HEADER   Add request header (repeatable)
//	-f          Follow redirects (default: false)
//	-k          Skip TLS verification
//	-n NAME     Print only specific header value
//	-j          JSON output
//	-t TIMEOUT  Request timeout (default: 10s)
//	-v          Verbose: also show request headers and status
//
// Examples:
//
//	headers https://example.com
//	headers -n content-type https://example.com
//	headers -f -j https://example.com
//	headers -m GET -H "Authorization: Bearer token" https://api.example.com
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type headerFlag []string

func (h *headerFlag) String() string  { return strings.Join(*h, ", ") }
func (h *headerFlag) Set(v string) error { *h = append(*h, v); return nil }

var extraHeaders headerFlag

func main() {
	method  := flag.String("m", "HEAD", "HTTP method")
	follow  := flag.Bool("f", false, "follow redirects")
	insecure := flag.Bool("k", false, "skip TLS verify")
	name    := flag.String("n", "", "specific header name")
	asJSON  := flag.Bool("j", false, "JSON output")
	timeout := flag.Duration("t", 10*time.Second, "timeout")
	verbose := flag.Bool("v", false, "verbose")
	flag.Var(&extraHeaders, "H", "request header")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: headers [OPTIONS] URL")
		os.Exit(1)
	}
	url := args[0]

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure},
	}
	client := &http.Client{
		Timeout:   *timeout,
		Transport: transport,
	}
	if !*follow {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	req, err := http.NewRequest(strings.ToUpper(*method), url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "headers: %v\n", err)
		os.Exit(1)
	}
	for _, h := range extraHeaders {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "headers: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if *name != "" {
		val := resp.Header.Get(*name)
		if val == "" {
			os.Exit(1)
		}
		fmt.Println(val)
		return
	}

	if *asJSON {
		m := make(map[string]interface{})
		m["status"] = resp.Status
		m["status_code"] = resp.StatusCode
		hdrs := make(map[string]string)
		for k, v := range resp.Header {
			hdrs[strings.ToLower(k)] = strings.Join(v, ", ")
		}
		m["headers"] = hdrs
		b, _ := json.MarshalIndent(m, "", "  ")
		fmt.Println(string(b))
		return
	}

	if *verbose {
		fmt.Printf("> %s %s\n", req.Method, url)
		for _, h := range extraHeaders {
			fmt.Printf("> %s\n", h)
		}
		fmt.Println()
	}

	fmt.Printf("< %s\n", resp.Status)
	// Sort header keys
	keys := make([]string, 0, len(resp.Header))
	for k := range resp.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range resp.Header[k] {
			fmt.Printf("< %s: %s\n", k, v)
		}
	}
}
