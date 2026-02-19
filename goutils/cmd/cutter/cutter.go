// cutter - cut bytes/chars/fields with more control than GNU cut
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func parseRanges(spec string) [][2]int {
	var ranges [][2]int
	for _, part := range strings.Split(spec, ",") {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "-"); idx >= 0 {
			lo, hi := part[:idx], part[idx+1:]
			a, _ := strconv.Atoi(lo)
			if hi == "" {
				ranges = append(ranges, [2]int{a, -1})
			} else {
				b, _ := strconv.Atoi(hi)
				ranges = append(ranges, [2]int{a, b})
			}
		} else {
			n, _ := strconv.Atoi(part)
			ranges = append(ranges, [2]int{n, n})
		}
	}
	sort.Slice(ranges, func(i, j int) bool { return ranges[i][0] < ranges[j][0] })
	return ranges
}

func inRanges(n int, ranges [][2]int) bool {
	for _, r := range ranges {
		if r[1] < 0 && n >= r[0] { return true }
		if n >= r[0] && n <= r[1] { return true }
	}
	return false
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: cutter [options] [file...]
  -b <list>   cut bytes (1-based, e.g. 1-5,10,20-)
  -c <list>   cut characters (unicode-aware)
  -f <list>   cut fields
  -d <delim>  field delimiter (default: tab)
  -o <delim>  output delimiter (default: same as input)
  -s          suppress lines with no delimiter (-f mode)
  --complement invert selection`)
	os.Exit(1)
}

func main() {
	mode := ""
	spec := ""
	delim := "\t"
	outDelim := ""
	suppress := false
	complement := false
	var files []string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-b": i++; mode = "b"; spec = args[i]
		case "-c": i++; mode = "c"; spec = args[i]
		case "-f": i++; mode = "f"; spec = args[i]
		case "-d": i++; delim = args[i]
		case "-o": i++; outDelim = args[i]
		case "-s": suppress = true
		case "--complement": complement = true
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			files = append(files, args[i])
		}
	}
	if mode == "" || spec == "" { usage() }
	if outDelim == "" { outDelim = delim }
	ranges := parseRanges(spec)

	process := func(line string) string {
		switch mode {
		case "b":
			bytes := []byte(line)
			var out []byte
			for i, b := range bytes {
				in := inRanges(i+1, ranges)
				if complement { in = !in }
				if in { out = append(out, b) }
			}
			return string(out)
		case "c":
			runes := []rune(line)
			var out []rune
			for i, r := range runes {
				in := inRanges(i+1, ranges)
				if complement { in = !in }
				if in { out = append(out, r) }
			}
			return string(out)
		case "f":
			if !strings.Contains(line, delim) {
				if suppress { return "" }
				return line
			}
			fields := strings.Split(line, delim)
			var out []string
			for i, f := range fields {
				in := inRanges(i+1, ranges)
				if complement { in = !in }
				if in { out = append(out, f) }
			}
			return strings.Join(out, outDelim)
		}
		return line
	}

	doFile := func(r *os.File) {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			out := process(sc.Text())
			if out != "" || !suppress { fmt.Println(out) }
		}
	}

	if len(files) == 0 { doFile(os.Stdin); return }
	for _, f := range files {
		fh, err := os.Open(f); if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		doFile(fh); fh.Close()
	}
}
