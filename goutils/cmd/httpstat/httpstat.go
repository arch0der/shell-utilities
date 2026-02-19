// httpstat - make HTTP requests and show detailed timing stats
package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"os"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: httpstat [options] <url>
  -X <method>       HTTP method (default: GET)
  -H <header>       add header (Key: Value), repeatable
  -d <body>         request body
  -k                skip TLS verification
  -f                follow redirects (default: yes)
  -t <timeout>      timeout duration (default: 30s)`)
	os.Exit(1)
}

func main() {
	method := "GET"
	var headers []string
	body := ""
	skipTLS := false
	timeout := 30 * time.Second
	var url string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-X": i++; method = strings.ToUpper(args[i])
		case "-H": i++; headers = append(headers, args[i])
		case "-d": i++; body = args[i]
		case "-k": skipTLS = true
		case "-t": i++; timeout, _ = time.ParseDuration(args[i])
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			url = args[i]
		}
	}
	if url == "" { usage() }
	if !strings.HasPrefix(url, "http") { url = "https://" + url }

	var t0, t1, t2, t3, t4, t5 time.Time
	trace := &httptrace.ClientTrace{
		DNSStart:          func(_ httptrace.DNSStartInfo) { t0 = time.Now() },
		DNSDone:           func(_ httptrace.DNSDoneInfo) { t1 = time.Now() },
		ConnectStart:      func(_, _ string) { t2 = time.Now() },
		ConnectDone:       func(_, _ string, _ error) { t3 = time.Now() },
		TLSHandshakeStart: func() { if t2.IsZero() { t2 = time.Now() } },
		TLSHandshakeDone:  func(_ tls.ConnectionState, _ error) { t4 = time.Now() },
		GotFirstResponseByte: func() { t5 = time.Now() },
	}

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLS},
		},
	}

	var bodyReader io.Reader
	if body != "" { bodyReader = strings.NewReader(body) }
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 { req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])) }
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	total := time.Since(start)

	fmt.Printf("\n%s %s\n\n", method, url)
	fmt.Printf("  Status    : %s\n", resp.Status)
	fmt.Printf("  Proto     : %s\n", resp.Proto)
	fmt.Printf("  Size      : %d bytes\n", len(respBody))
	fmt.Println()
	if !t0.IsZero() && !t1.IsZero() { fmt.Printf("  DNS Lookup    : %v\n", t1.Sub(t0)) }
	if !t2.IsZero() && !t3.IsZero() { fmt.Printf("  TCP Connect   : %v\n", t3.Sub(t2)) }
	if !t3.IsZero() && !t4.IsZero() { fmt.Printf("  TLS Handshake : %v\n", t4.Sub(t3)) }
	if !t5.IsZero() { fmt.Printf("  TTFB          : %v\n", t5.Sub(start)) }
	fmt.Printf("  Total         : %v\n\n", total)
	fmt.Println("  Response Headers:")
	for k, vs := range resp.Header { fmt.Printf("    %-30s %s\n", k+":", strings.Join(vs, ", ")) }
}
