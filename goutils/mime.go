// mime - Detect MIME type of files.
//
// Usage:
//
//	mime [OPTIONS] FILE [FILE...]
//	cat file | mime
//
// Options:
//
//	-e        Print file extension for the detected type
//	-b        Brief: print MIME type only, no filename
//	-j        JSON output
//	-0        Exit 0 even on unknown types
//
// Examples:
//
//	mime image.png              # image/png
//	mime -e image.png           # .png
//	mime -b *.jpg               # image/jpeg (for each)
//	cat unknown | mime          # detect from stdin
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	ext    = flag.Bool("e", false, "print extension")
	brief  = flag.Bool("b", false, "brief output")
	asJSON = flag.Bool("j", false, "JSON output")
)

// Extended MIME map for common types that net/http might miss
var extMIME = map[string]string{
	".jpg": "image/jpeg", ".jpeg": "image/jpeg",
	".png": "image/png", ".gif": "image/gif",
	".webp": "image/webp", ".svg": "image/svg+xml",
	".ico": "image/x-icon", ".bmp": "image/bmp",
	".mp4": "video/mp4", ".mkv": "video/x-matroska",
	".mov": "video/quicktime", ".avi": "video/x-msvideo",
	".mp3": "audio/mpeg", ".wav": "audio/wav",
	".ogg": "audio/ogg", ".flac": "audio/flac",
	".pdf": "application/pdf",
	".zip": "application/zip", ".tar": "application/x-tar",
	".gz": "application/gzip", ".bz2": "application/x-bzip2",
	".xz": "application/x-xz", ".7z": "application/x-7z-compressed",
	".json": "application/json", ".xml": "application/xml",
	".yaml": "application/yaml", ".yml": "application/yaml",
	".toml": "application/toml", ".csv": "text/csv",
	".html": "text/html", ".htm": "text/html",
	".css": "text/css", ".js": "text/javascript",
	".ts": "text/typescript", ".go": "text/x-go",
	".py": "text/x-python", ".rb": "text/x-ruby",
	".sh": "text/x-shellscript", ".md": "text/markdown",
	".txt": "text/plain", ".log": "text/plain",
	".exe": "application/vnd.microsoft.portable-executable",
	".dll": "application/vnd.microsoft.portable-executable",
	".so": "application/x-sharedlib",
	".wasm": "application/wasm",
	".db": "application/x-sqlite3", ".sqlite": "application/x-sqlite3",
}

var mimeExt = map[string]string{}

func init() {
	for e, m := range extMIME {
		if _, ok := mimeExt[m]; !ok {
			mimeExt[m] = e
		}
	}
}

func detectMIME(path string) string {
	// try by extension first
	e := strings.ToLower(filepath.Ext(path))
	if m, ok := extMIME[e]; ok {
		return m
	}
	// fall back to reading file magic bytes
	f, err := os.Open(path)
	if err != nil {
		return "application/octet-stream"
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	return http.DetectContentType(buf[:n])
}

type Result struct {
	File string `json:"file"`
	MIME string `json:"mime"`
	Ext  string `json:"ext,omitempty"`
}

func main() {
	flag.Parse()
	files := flag.Args()

	process := func(path string) Result {
		m := detectMIME(path)
		e := mimeExt[m]
		return Result{File: path, MIME: m, Ext: e}
	}

	if len(files) == 0 {
		// stdin
		buf := make([]byte, 512)
		n, _ := os.Stdin.Read(buf)
		m := http.DetectContentType(buf[:n])
		if *brief {
			fmt.Println(m)
			return
		}
		fmt.Printf("stdin: %s\n", m)
		return
	}

	var results []Result
	for _, f := range files {
		r := process(f)
		results = append(results, r)
	}

	if *asJSON {
		b, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(b))
		return
	}

	for _, r := range results {
		if *ext {
			if *brief {
				fmt.Println(r.Ext)
			} else {
				fmt.Printf("%s: %s\n", r.File, r.Ext)
			}
			continue
		}
		if *brief {
			fmt.Println(r.MIME)
		} else {
			fmt.Printf("%s: %s\n", r.File, r.MIME)
		}
	}
}
