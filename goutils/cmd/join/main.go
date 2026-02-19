// join - Join lines of two files on a common field
// Usage: join [-1 field] [-2 field] [-t sep] [-a 1|2] file1 file2
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	field1 = flag.Int("1", 1, "Join field in file1 (1-indexed)")
	field2 = flag.Int("2", 1, "Join field in file2 (1-indexed)")
	sep    = flag.String("t", " ", "Field separator")
	unpaired = flag.Int("a", 0, "Print unpairable lines from file N (1 or 2)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: join [-1 field] [-2 field] [-t sep] [-a 1|2] file1 file2")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	lines1 := readIndexed(flag.Arg(0), *field1-1, *sep)
	lines2 := readIndexed(flag.Arg(1), *field2-1, *sep)

	matched2 := map[string]bool{}
	for key, rows1 := range lines1 {
		rows2, ok := lines2[key]
		if ok {
			for _, r1 := range rows1 {
				for _, r2 := range rows2 {
					// Merge: key + rest of r1 + rest of r2
					out := []string{key}
					out = append(out, withoutField(r1, *field1-1)...)
					out = append(out, withoutField(r2, *field2-1)...)
					fmt.Println(strings.Join(out, *sep))
				}
			}
			matched2[key] = true
		} else if *unpaired == 1 {
			for _, r1 := range rows1 {
				fmt.Println(strings.Join(r1, *sep))
			}
		}
	}
	if *unpaired == 2 {
		for key, rows2 := range lines2 {
			if !matched2[key] {
				for _, r2 := range rows2 {
					fmt.Println(strings.Join(r2, *sep))
				}
			}
		}
	}
}

func readIndexed(path string, fieldIdx int, sep string) map[string][][]string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "join:", err)
		os.Exit(1)
	}
	defer f.Close()
	result := map[string][][]string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Split(sc.Text(), sep)
		if fieldIdx >= len(fields) {
			continue
		}
		key := fields[fieldIdx]
		result[key] = append(result[key], fields)
	}
	return result
}

func withoutField(fields []string, idx int) []string {
	out := []string{}
	for i, f := range fields {
		if i != idx {
			out = append(out, f)
		}
	}
	return out
}
