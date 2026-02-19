// bsearch - binary search in a sorted file or stdin for a key
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	caseSensitive := true
	fieldN := 0 // 0 = whole line, else 1-based field index
	delim := "\t"
	args := os.Args[1:]
	rest := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-i": caseSensitive = false
		case "-f": i++; fmt.Sscanf(args[i], "%d", &fieldN)
		case "-d": i++; delim = args[i]
		default: rest = append(rest, args[i])
		}
	}
	if len(rest) < 1 {
		fmt.Fprintln(os.Stderr, "usage: bsearch [options] <key> [sorted_file]")
		fmt.Fprintln(os.Stderr, "  -i      case-insensitive")
		fmt.Fprintln(os.Stderr, "  -f <n>  match field n (1-based, tab-separated)")
		fmt.Fprintln(os.Stderr, "  -d <d>  field delimiter (default: tab)")
		os.Exit(1)
	}
	key := rest[0]
	if !caseSensitive { key = strings.ToLower(key) }

	var r *os.File
	if len(rest) > 1 {
		var err error; r, err = os.Open(rest[1]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		defer r.Close()
	} else { r = os.Stdin }

	sc := bufio.NewScanner(r)
	var lines []string
	for sc.Scan() { lines = append(lines, sc.Text()) }

	getKey := func(line string) string {
		k := line
		if fieldN > 0 {
			parts := strings.Split(line, delim)
			if fieldN <= len(parts) { k = parts[fieldN-1] }
		}
		if !caseSensitive { k = strings.ToLower(k) }
		return k
	}

	idx := sort.Search(len(lines), func(i int) bool { return getKey(lines[i]) >= key })
	found := false
	for i := idx; i < len(lines) && getKey(lines[i]) == key; i++ {
		fmt.Println(lines[i]); found = true
	}
	if !found { fmt.Fprintf(os.Stderr, "bsearch: %q not found\n", key); os.Exit(1) }
}
