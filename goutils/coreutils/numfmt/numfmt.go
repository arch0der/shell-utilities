package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	fromUnit := ""
	toUnit := ""
	suffix := ""
	padding := 0
	headerLines := 0
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case strings.HasPrefix(a, "--from="):
			fromUnit = a[7:]
		case strings.HasPrefix(a, "--to="):
			toUnit = a[5:]
		case strings.HasPrefix(a, "--suffix="):
			suffix = a[9:]
		case strings.HasPrefix(a, "--padding="):
			fmt.Sscan(a[10:], &padding)
		case strings.HasPrefix(a, "--header="):
			fmt.Sscan(a[9:], &headerLines)
		case a == "--header":
			headerLines = 1
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}

	parseNum := func(s string) float64 {
		s = strings.TrimSpace(s)
		suffixes := map[byte]float64{
			'K': 1e3, 'M': 1e6, 'G': 1e9, 'T': 1e12, 'P': 1e15, 'E': 1e18,
		}
		iSuffixes := map[byte]float64{
			'K': 1024, 'M': 1024 * 1024, 'G': 1024 * 1024 * 1024,
		}
		multiplier := 1.0
		if fromUnit == "iec" || fromUnit == "iec-i" || fromUnit == "auto" {
			if len(s) > 1 {
				last := s[len(s)-1]
				if mult, ok := iSuffixes[last]; ok {
					multiplier = mult
					s = s[:len(s)-1]
				}
			}
		} else if len(s) > 1 {
			last := s[len(s)-1]
			if mult, ok := suffixes[last]; ok {
				multiplier = mult
				s = s[:len(s)-1]
			}
		}
		n, _ := strconv.ParseFloat(s, 64)
		return n * multiplier
	}

	formatNum := func(n float64) string {
		if toUnit == "" {
			return strconv.FormatFloat(n, 'f', -1, 64)
		}
		units := []struct {
			suffix string
			value  float64
		}{
			{"E", 1e18}, {"P", 1e15}, {"T", 1e12}, {"G", 1e9}, {"M", 1e6}, {"K", 1e3}, {"", 1},
		}
		if toUnit == "iec" || toUnit == "iec-i" {
			units = []struct {
				suffix string
				value  float64
			}{
				{"E", math.Pow(1024, 6)}, {"P", math.Pow(1024, 5)},
				{"T", math.Pow(1024, 4)}, {"G", math.Pow(1024, 3)},
				{"M", math.Pow(1024, 2)}, {"K", 1024}, {"", 1},
			}
		}
		for _, u := range units {
			if math.Abs(n) >= u.value {
				val := n / u.value
				if val == math.Trunc(val) {
					return fmt.Sprintf("%.0f%s%s", val, u.suffix, suffix)
				}
				return fmt.Sprintf("%.1f%s%s", val, u.suffix, suffix)
			}
		}
		return fmt.Sprintf("%.0f%s", n, suffix)
	}

	processLine := func(line string) string {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			return line
		}
		n := parseNum(fields[0])
		result := formatNum(n)
		if padding != 0 {
			result = fmt.Sprintf("%*s", padding, result)
		}
		if len(fields) > 1 {
			return result + " " + strings.Join(fields[1:], " ")
		}
		return result
	}

	process := func(r io.Reader) {
		sc := bufio.NewScanner(r)
		headerCount := 0
		for sc.Scan() {
			line := sc.Text()
			if headerCount < headerLines {
				fmt.Println(line)
				headerCount++
				continue
			}
			fmt.Println(processLine(line))
		}
	}

	if len(files) == 0 {
		process(os.Stdin)
		return
	}
	for _, f := range files {
		fh, _ := os.Open(f)
		process(fh)
		fh.Close()
	}
}
