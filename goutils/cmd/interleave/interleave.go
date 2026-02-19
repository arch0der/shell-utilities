// interleave - interleave lines from multiple files (round-robin or by ratio)
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: interleave [options] <file1> <file2> [file3...]
  -n <n>    lines from each file per turn (default: 1)
  -r <a:b>  ratio of lines (for 2 files only), e.g. 2:1
  -s <sep>  separator line between groups (default: none)
  -p        pad with empty lines to keep alignment`)
	os.Exit(1)
}

func main() {
	n := 1
	ratio := ""
	sep := ""
	pad := false
	var files []string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-n": i++; n, _ = strconv.Atoi(args[i])
		case "-r": i++; ratio = args[i]
		case "-s": i++; sep = args[i]
		case "-p": pad = true
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			files = append(files, args[i])
		}
	}
	if len(files) < 2 { usage() }

	// Read all files
	allLines := make([][]string, len(files))
	for i, f := range files {
		fh, err := os.Open(f); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		sc := bufio.NewScanner(fh)
		for sc.Scan() { allLines[i] = append(allLines[i], sc.Text()) }
		fh.Close()
	}

	// Compute ratios
	ns := make([]int, len(files))
	for i := range ns { ns[i] = n }
	if ratio != "" && len(files) == 2 {
		parts := strings.Split(ratio, ":")
		if len(parts) == 2 { ns[0], _ = strconv.Atoi(parts[0]); ns[1], _ = strconv.Atoi(parts[1]) }
	}

	// Interleave
	indices := make([]int, len(files))
	maxLen := 0; for _, l := range allLines { if len(l) > maxLen { maxLen = len(l) } }
	printed := false
	for {
		anyLeft := false
		for fi, lines := range allLines {
			start := indices[fi]
			for j := 0; j < ns[fi]; j++ {
				idx := start + j
				if idx < len(lines) { fmt.Println(lines[idx]); anyLeft = true } else if pad { fmt.Println() }
			}
			indices[fi] += ns[fi]
			if sep != "" && fi < len(allLines)-1 && anyLeft { fmt.Println(sep) }
		}
		if sep != "" && printed { fmt.Println(sep) }
		_ = printed; printed = true
		// Check if all done
		done := true
		for fi, lines := range allLines { if indices[fi] < len(lines) { done = false; break } }
		if done { break }
	}
}
