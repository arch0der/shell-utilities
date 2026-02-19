// csplit - Split a file into sections by context (pattern)
// Usage: csplit [-f prefix] [-n digits] [-k] file pattern [pattern...]
// Patterns: /regex/ - split before match, %regex% - skip to match, N - line number
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	prefix    = flag.String("f", "xx", "Output file prefix")
	digits    = flag.Int("n", 2, "Number of digits in output filenames")
	keep      = flag.Bool("k", false, "Keep output files on error")
	quiet     = flag.Bool("q", false, "Suppress byte counts")
	elide     = flag.Bool("z", false, "Remove empty output files")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: csplit [-f prefix] [-n digits] [-k] [-q] [-z] file pattern...")
		fmt.Fprintln(os.Stderr, "Patterns: /regex/ split before match, N line number")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	patterns := flag.Args()[1:]

	var lines []string
	if filename == "-" {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
	} else {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "csplit:", err)
			os.Exit(1)
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			lines = append(lines, sc.Text())
		}
		f.Close()
	}

	// Find split points
	splitPoints := []int{0}
	for _, pat := range patterns {
		if strings.HasPrefix(pat, "/") && strings.HasSuffix(pat, "/") {
			re := regexp.MustCompile(pat[1 : len(pat)-1])
			// Find next match after last split point
			last := splitPoints[len(splitPoints)-1]
			for i := last; i < len(lines); i++ {
				if re.MatchString(lines[i]) {
					splitPoints = append(splitPoints, i)
					break
				}
			}
		} else if n, err := strconv.Atoi(pat); err == nil {
			splitPoints = append(splitPoints, n-1)
		}
	}
	splitPoints = append(splitPoints, len(lines))

	// Write sections
	fmtStr := fmt.Sprintf("%%s%%0%dd", *digits)
	created := []string{}
	for i := 0; i < len(splitPoints)-1; i++ {
		start := splitPoints[i]
		end := splitPoints[i+1]
		section := lines[start:end]

		if *elide && len(section) == 0 {
			continue
		}

		name := fmt.Sprintf(fmtStr, *prefix, i)
		created = append(created, name)

		content := strings.Join(section, "\n")
		if len(section) > 0 {
			content += "\n"
		}

		if err := os.WriteFile(name, []byte(content), 0644); err != nil {
			fmt.Fprintln(os.Stderr, "csplit:", err)
			if !*keep {
				for _, c := range created {
					os.Remove(c)
				}
				os.Exit(1)
			}
		}

		if !*quiet {
			fmt.Println(len(content))
		}
	}
}
