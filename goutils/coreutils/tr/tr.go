package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func expandTrSet(s string) []rune {
	var result []rune
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if i+2 < len(runes) && runes[i+1] == '-' {
			for c := runes[i]; c <= runes[i+2]; c++ {
				result = append(result, c)
			}
			i += 2
		} else if runes[i] == '\\' && i+1 < len(runes) {
			i++
			switch runes[i] {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '\\':
				result = append(result, '\\')
			default:
				result = append(result, runes[i])
			}
		} else {
			result = append(result, runes[i])
		}
	}
	return result
}

func main() {
	args := os.Args[1:]
	deleteMode := false
	squeezeMode := false
	complement := false
	sets := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-d" || a == "--delete":
			deleteMode = true
		case a == "-s" || a == "--squeeze-repeats":
			squeezeMode = true
		case a == "-c" || a == "-C" || a == "--complement":
			complement = true
		case !strings.HasPrefix(a, "-"):
			sets = append(sets, a)
		}
	}

	set1 := []rune{}
	set2 := []rune{}
	if len(sets) > 0 {
		set1 = expandTrSet(sets[0])
	}
	if len(sets) > 1 {
		set2 = expandTrSet(sets[1])
	}

	// Build complement of set1
	if complement && !deleteMode {
		allChars := make([]rune, 256)
		for i := range allChars {
			allChars[i] = rune(i)
		}
		set1Map := map[rune]bool{}
		for _, c := range set1 {
			set1Map[c] = true
		}
		var comp []rune
		for _, c := range allChars {
			if !set1Map[c] {
				comp = append(comp, c)
			}
		}
		set1 = comp
	}

	inSet1 := func(c rune) (int, bool) {
		for i, r := range set1 {
			if r == c {
				return i, true
			}
		}
		return -1, false
	}

	data, _ := io.ReadAll(os.Stdin)
	var out strings.Builder
	prev := rune(-1)

	for _, c := range string(data) {
		if deleteMode {
			idx, found := inSet1(c)
			_ = idx
			if complement {
				if !found {
					continue
				}
			} else {
				if found {
					continue
				}
			}
			out.WriteRune(c)
			continue
		}

		replacement := c
		if idx, found := inSet1(c); found && len(set2) > 0 {
			si := idx
			if si >= len(set2) {
				si = len(set2) - 1
			}
			replacement = set2[si]
		}

		if squeezeMode && replacement == prev {
			continue
		}
		out.WriteRune(replacement)
		prev = replacement
	}
	fmt.Print(out.String())
}
