// numfmt - Format numbers with commas, SI suffixes, or byte units.
//
// Usage:
//
//	numfmt [OPTIONS] [NUMBER...]
//	echo "1234567" | numfmt
//
// Options:
//
//	-s        Add thousands separator (default)
//	-si       Format as SI units (K, M, G, T)
//	-iec      Format as IEC units (Ki, Mi, Gi, Ti)
//	-from-si  Parse SI suffix back to raw number
//	-d SEP    Input delimiter for multi-column (default: whitespace)
//	-f N      Format field N (1-based, default: last field)
//	-p N      Decimal precision (default: 2 for SI/IEC)
//
// Examples:
//
//	echo 1234567 | numfmt               # 1,234,567
//	echo 1073741824 | numfmt -iec       # 1.00Gi
//	echo 1500000 | numfmt -si           # 1.50M
//	echo "1.5K" | numfmt -from-si       # 1500
//	df -B1 | numfmt -f 2 -iec           # format 2nd field as IEC
package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	useSI    = flag.Bool("si", false, "SI units")
	useIEC   = flag.Bool("iec", false, "IEC units")
	fromSI   = flag.Bool("from-si", false, "parse SI suffix")
	delim    = flag.String("d", "", "delimiter")
	field    = flag.Int("f", 0, "field to format (1-based)")
	prec     = flag.Int("p", 2, "decimal precision")
)

func addCommas(s string) string {
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}
	dot := strings.Index(s, ".")
	intPart, fracPart := s, ""
	if dot >= 0 {
		intPart, fracPart = s[:dot], s[dot:]
	}
	var out []byte
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	result := string(out) + fracPart
	if neg {
		return "-" + result
	}
	return result
}

func formatSI(n float64, iec bool) string {
	suffixes := []string{"", "K", "M", "G", "T", "P", "E"}
	iecSuf := []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei"}
	base := 1000.0
	suf := suffixes
	if iec {
		base = 1024.0
		suf = iecSuf
	}
	if math.Abs(n) < base {
		return fmt.Sprintf("%d", int64(n))
	}
	div := n
	i := 0
	for math.Abs(div) >= base && i < len(suf)-1 {
		div /= base
		i++
	}
	return fmt.Sprintf("%.*f%s", *prec, div, suf[i])
}

func parseSI(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty")
	}
	mults := map[byte]float64{
		'K': 1e3, 'M': 1e6, 'G': 1e9, 'T': 1e12, 'P': 1e15,
		'k': 1e3, 'm': 1e6, 'g': 1e9, 't': 1e12, 'p': 1e15,
	}
	last := s[len(s)-1]
	if mult, ok := mults[last]; ok {
		n, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, err
		}
		return n * mult, nil
	}
	return strconv.ParseFloat(s, 64)
}

func process(tok string) string {
	if *fromSI {
		n, err := parseSI(tok)
		if err != nil {
			return tok
		}
		return fmt.Sprintf("%g", n)
	}
	n, err := strconv.ParseFloat(strings.ReplaceAll(tok, ",", ""), 64)
	if err != nil {
		return tok
	}
	if *useSI {
		return formatSI(n, false)
	}
	if *useIEC {
		return formatSI(n, true)
	}
	// default: commas
	s := fmt.Sprintf("%.0f", n)
	return addCommas(s)
}

func main() {
	flag.Parse()
	args := flag.Args()

	processLine := func(line string) string {
		sep := " "
		if *delim != "" {
			sep = *delim
		}
		parts := strings.Split(line, sep)
		if *field > 0 && *field <= len(parts) {
			parts[*field-1] = process(parts[*field-1])
			return strings.Join(parts, sep)
		}
		// no field spec: process whole line as number
		return process(strings.TrimSpace(line))
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if len(args) > 0 {
		for _, a := range args {
			fmt.Fprintln(w, process(a))
		}
		return
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Fprintln(w, processLine(sc.Text()))
	}
}
