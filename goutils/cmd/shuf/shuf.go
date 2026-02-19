// shuf - Output lines in random order (shuffle).
//
// Usage:
//
//	shuf [OPTIONS] [FILE...]
//	echo -e "a\nb\nc" | shuf
//
// Options:
//
//	-n N      Output at most N lines
//	-r        Allow repetition (sample with replacement)
//	-i LO-HI  Shuffle integers in range LO-HI instead of lines
//	-e ARGS   Treat each argument as a line
//	-z        Output null-separated instead of newlines
//	--seed N  Random seed (for reproducibility)
//
// Examples:
//
//	cat list.txt | shuf                 # shuffle all lines
//	shuf -n 3 words.txt                 # pick 3 random lines
//	shuf -i 1-100 -n 5                  # 5 random numbers 1-100
//	shuf -e apple banana cherry         # shuffle args
//	shuf -r -n 10 colors.txt           # 10 with repetition
package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	count    = flag.Int("n", 0, "max output lines (0=all)")
	repeat   = flag.Bool("r", false, "allow repetition")
	intRange = flag.String("i", "", "integer range LO-HI")
	echo     = flag.Bool("e", false, "treat args as lines")
	null     = flag.Bool("z", false, "null-separated output")
	seed     = flag.Int64("seed", 0, "random seed")
)

func main() {
	flag.Parse()

	src := rand.NewSource(time.Now().UnixNano())
	if *seed != 0 {
		src = rand.NewSource(*seed)
	}
	rng := rand.New(src)

	var lines []string

	// Integer range mode
	if *intRange != "" {
		parts := strings.SplitN(*intRange, "-", 2)
		if len(parts) != 2 {
			fmt.Fprintln(os.Stderr, "shuf: invalid range, use LO-HI")
			os.Exit(1)
		}
		lo, _ := strconv.Atoi(parts[0])
		hi, _ := strconv.Atoi(parts[1])
		for i := lo; i <= hi; i++ {
			lines = append(lines, strconv.Itoa(i))
		}
	} else if *echo {
		lines = flag.Args()
	} else {
		files := flag.Args()
		var readers []*os.File
		if len(files) == 0 {
			readers = []*os.File{os.Stdin}
		} else {
			for _, f := range files {
				fh, err := os.Open(f)
				if err != nil {
					fmt.Fprintf(os.Stderr, "shuf: %v\n", err)
					os.Exit(1)
				}
				defer fh.Close()
				readers = append(readers, fh)
			}
		}
		for _, r := range readers {
			sc := bufio.NewScanner(r)
			for sc.Scan() {
				lines = append(lines, sc.Text())
			}
		}
	}

	if len(lines) == 0 {
		return
	}

	sep := "\n"
	if *null {
		sep = "\x00"
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	n := *count
	if n == 0 || n > len(lines) {
		n = len(lines)
	}

	if *repeat {
		for i := 0; i < n; i++ {
			idx := rng.Intn(len(lines))
			fmt.Fprintf(w, "%s%s", lines[idx], sep)
		}
		return
	}

	// Fisher-Yates shuffle
	perm := rng.Perm(len(lines))
	for i := 0; i < n; i++ {
		fmt.Fprintf(w, "%s%s", lines[perm[i]], sep)
	}
}
