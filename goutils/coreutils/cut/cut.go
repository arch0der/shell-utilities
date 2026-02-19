package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type cutRange struct{ lo, hi int }

func parseCutRanges(s string) []cutRange {
	var ranges []cutRange
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			lr := strings.SplitN(part, "-", 2)
			lo, hi := 0, 0
			if lr[0] != "" {
				lo, _ = strconv.Atoi(lr[0])
			} else {
				lo = 1
			}
			if lr[1] != "" {
				hi, _ = strconv.Atoi(lr[1])
			}
			ranges = append(ranges, cutRange{lo, hi})
		} else {
			n, _ := strconv.Atoi(part)
			ranges = append(ranges, cutRange{n, n})
		}
	}
	return ranges
}

func main() {
	args := os.Args[1:]
	delim := "\t"
	var fieldSpec string
	var byteSpec string
	var charSpec string
	complement := false
	onlyDelim := false
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-d" && i+1 < len(args):
			i++
			delim = args[i]
		case strings.HasPrefix(a, "-d"):
			delim = a[2:]
		case a == "-f" && i+1 < len(args):
			i++
			fieldSpec = args[i]
		case strings.HasPrefix(a, "-f"):
			fieldSpec = a[2:]
		case a == "-b" && i+1 < len(args):
			i++
			byteSpec = args[i]
		case strings.HasPrefix(a, "-b"):
			byteSpec = a[2:]
		case a == "-c" && i+1 < len(args):
			i++
			charSpec = args[i]
		case strings.HasPrefix(a, "-c"):
			charSpec = a[2:]
		case a == "--complement":
			complement = true
		case a == "-s" || a == "--only-delimited":
			onlyDelim = true
		case !strings.HasPrefix(a, "-") || a == "-":
			files = append(files, a)
		}
	}

	processLine := func(line string) string {
		if fieldSpec != "" {
			if onlyDelim && !strings.Contains(line, delim) {
				return ""
			}
			parts := strings.Split(line, delim)
			ranges := parseCutRanges(fieldSpec)
			inRange := func(n int) bool {
				for _, r := range ranges {
					hi := r.hi
					if hi == 0 {
						hi = len(parts)
					}
					if n >= r.lo && n <= hi {
						return true
					}
				}
				return false
			}
			var selected []string
			for i, p := range parts {
				inc := inRange(i + 1)
				if complement {
					inc = !inc
				}
				if inc {
					selected = append(selected, p)
				}
			}
			return strings.Join(selected, delim)
		}
		// byte or char mode
		spec := byteSpec
		if spec == "" {
			spec = charSpec
		}
		runes := []rune(line)
		ranges := parseCutRanges(spec)
		// collect indices
		inRange := func(n int) bool {
			for _, r := range ranges {
				hi := r.hi
				if hi == 0 {
					hi = len(runes)
				}
				if n >= r.lo && n <= hi {
					return true
				}
			}
			return false
		}
		// collect unique sorted indices
		idxSet := map[int]bool{}
		for _, r := range ranges {
			lo, hi := r.lo, r.hi
			if hi == 0 {
				hi = len(runes)
			}
			for i := lo; i <= hi; i++ {
				idxSet[i] = true
			}
		}
		idxList := make([]int, 0, len(idxSet))
		for k := range idxSet {
			idxList = append(idxList, k)
		}
		sort.Ints(idxList)
		var sb strings.Builder
		for _, idx := range idxList {
			inc := inRange(idx)
			if complement {
				inc = !inc
			}
			if inc && idx >= 1 && idx <= len(runes) {
				sb.WriteRune(runes[idx-1])
			}
		}
		return sb.String()
	}

	processReader := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for sc.Scan() {
			result := processLine(sc.Text())
			if result != "" || fieldSpec == "" {
				fmt.Println(result)
			}
		}
	}

	if len(files) == 0 {
		processReader(os.Stdin)
		return
	}
	for _, f := range files {
		if f == "-" {
			processReader(os.Stdin)
			continue
		}
		fh, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cut: %s: %v\n", f, err)
			continue
		}
		processReader(fh)
		fh.Close()
	}
}
