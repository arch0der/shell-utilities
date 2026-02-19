// floatfmt - format floating point numbers with controlled precision and style
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: floatfmt [options] [numbers...]
  -p <n>     decimal precision (default: 2)
  -s         use scientific notation
  -c         add thousands comma separator
  -pct       multiply by 100 and add % sign
  -trim      trim trailing zeros`)
	os.Exit(1)
}

func addCommas(s string) string {
	parts := strings.SplitN(s, ".", 2)
	n := parts[0]
	neg := strings.HasPrefix(n, "-")
	if neg { n = n[1:] }
	var result []byte
	for i, ch := range []byte(n) {
		if i > 0 && (len(n)-i)%3 == 0 { result = append(result, ',') }
		result = append(result, ch)
	}
	out := string(result)
	if neg { out = "-" + out }
	if len(parts) > 1 { out += "." + parts[1] }
	return out
}

func main() {
	prec := 2
	sci := false
	commas := false
	pct := false
	trim := false
	var nums []string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-p": i++; prec, _ = strconv.Atoi(args[i])
		case "-s": sci = true
		case "-c": commas = true
		case "-pct": pct = true
		case "-trim": trim = true
		default:
			if strings.HasPrefix(args[i], "-") { usage() }
			nums = append(nums, args[i])
		}
	}

	format := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" { return }
		f, err := strconv.ParseFloat(s, 64)
		if err != nil { fmt.Fprintf(os.Stderr, "floatfmt: %q is not a number\n", s); return }
		if pct { f *= 100 }
		var out string
		if sci {
			out = fmt.Sprintf("%.*e", prec, f)
		} else {
			out = fmt.Sprintf("%.*f", prec, f)
		}
		if trim {
			if strings.Contains(out, ".") {
				out = strings.TrimRight(out, "0")
				out = strings.TrimRight(out, ".")
			}
		}
		if commas && !sci { out = addCommas(out) }
		if pct { out += "%" }
		fmt.Println(out)
	}

	if len(nums) > 0 {
		for _, n := range nums { format(n) }
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { format(sc.Text()) }
}
