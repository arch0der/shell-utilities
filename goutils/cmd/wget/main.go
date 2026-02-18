// wget - Download files from URLs
// Usage: wget [-O output] [-q] <url>
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	output = flag.String("O", "", "Output file name")
	quiet  = flag.Bool("q", false, "Quiet mode")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: wget [-O output] [-q] <url>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	url := flag.Arg(0)
	outFile := *output
	if outFile == "" {
		// Derive filename from URL
		outFile = path.Base(url)
		if outFile == "." || outFile == "/" {
			outFile = "index.html"
		}
		// Strip query string
		if i := strings.Index(outFile, "?"); i >= 0 {
			outFile = outFile[:i]
		}
	}

	if !*quiet {
		fmt.Fprintf(os.Stderr, "---> %s\n", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wget:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "wget: HTTP %d\n", resp.StatusCode)
		os.Exit(1)
	}

	f, err := os.Create(outFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wget:", err)
		os.Exit(1)
	}
	defer f.Close()

	written, err := io.Copy(f, resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wget:", err)
		os.Exit(1)
	}

	if !*quiet {
		fmt.Fprintf(os.Stderr, "Saved '%s' (%d bytes)\n", outFile, written)
	}
}
