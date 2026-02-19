// crossref - build a cross-reference index: show which lines each word appears on
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	minLen := 3
	ignoreCase := false
	var stopWords = map[string]bool{
		"the":true,"and":true,"for":true,"are":true,"but":true,"not":true,
		"you":true,"all":true,"can":true,"had":true,"her":true,"was":true,
		"one":true,"our":true,"out":true,"day":true,"get":true,"has":true,
		"him":true,"his":true,"how":true,"its":true,"may":true,"now":true,
		"did":true,"its":true,"let":true,"man":true,"new":true,"old":true,
		"see":true,"two":true,"way":true,"who":true,"boy":true,"did":true,
		"this":true,"that":true,"with":true,"have":true,"from":true,"they":true,
		"will":true,"been":true,"into":true,"more":true,"were":true,"what":true,
		"when":true,"your":true,"said":true,"each":true,"which":true,"she":true,
		"then":true,"than":true,"them":true,"these":true,"there":true,
	}
	files := []string{}
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-i": ignoreCase = true
		case "-m": i++; fmt.Sscanf(args[i], "%d", &minLen)
		case "--no-stop": stopWords = map[string]bool{}
		default: files = append(files, args[i])
		}
	}

	wordRe := regexp.MustCompile(`[a-zA-Z']+`)
	index := map[string][]int{}
	lineNum := 1

	process := func(r *os.File) {
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			line := sc.Text()
			words := wordRe.FindAllString(line, -1)
			seen := map[string]bool{}
			for _, w := range words {
				if ignoreCase { w = strings.ToLower(w) }
				w = strings.Trim(w, "'")
				if len(w) < minLen || stopWords[strings.ToLower(w)] || seen[w] { continue }
				seen[w] = true
				index[w] = append(index[w], lineNum)
			}
			lineNum++
		}
	}

	if len(files) == 0 { process(os.Stdin) } else {
		for _, f := range files {
			fh, err := os.Open(f); if err != nil { fmt.Fprintln(os.Stderr, err); continue }
			process(fh); fh.Close()
		}
	}

	keys := make([]string, 0, len(index))
	for k := range index { keys = append(keys, k) }
	sort.Strings(keys)
	for _, k := range keys {
		lines := index[k]
		strs := make([]string, len(lines))
		for i, l := range lines { strs[i] = fmt.Sprintf("%d", l) }
		fmt.Printf("%-24s %s\n", k, strings.Join(strs, ", "))
	}
}
