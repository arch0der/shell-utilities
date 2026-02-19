// roman - convert between integers and Roman numerals
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var toRomanTable = []struct{ v int; s string }{
	{1000,"M"},{900,"CM"},{500,"D"},{400,"CD"},
	{100,"C"},{90,"XC"},{50,"L"},{40,"XL"},
	{10,"X"},{9,"IX"},{5,"V"},{4,"IV"},{1,"I"},
}

func toRoman(n int) (string, error) {
	if n < 1 || n > 3999 { return "", fmt.Errorf("out of range (1-3999): %d", n) }
	var sb strings.Builder
	for _, t := range toRomanTable { for n >= t.v { sb.WriteString(t.s); n -= t.v } }
	return sb.String(), nil
}

var fromRomanMap = map[byte]int{'I':1,'V':5,'X':10,'L':50,'C':100,'D':500,'M':1000}

func fromRoman(s string) (int, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" { return 0, fmt.Errorf("empty input") }
	total, prev := 0, 0
	for i := len(s) - 1; i >= 0; i-- {
		v, ok := fromRomanMap[s[i]]
		if !ok { return 0, fmt.Errorf("invalid character %q", s[i]) }
		if v < prev { total -= v } else { total += v }
		prev = v
	}
	if total < 1 || total > 3999 { return 0, fmt.Errorf("result out of range") }
	return total, nil
}

func process(s string) {
	s = strings.TrimSpace(s)
	if s == "" { return }
	if n, err := strconv.Atoi(s); err == nil {
		r, err := toRoman(n)
		if err != nil { fmt.Fprintln(os.Stderr, "roman:", err); return }
		fmt.Printf("%d → %s\n", n, r)
	} else {
		n, err := fromRoman(s)
		if err != nil { fmt.Fprintln(os.Stderr, "roman:", err); return }
		fmt.Printf("%s → %d\n", strings.ToUpper(s), n)
	}
}

func main() {
	if len(os.Args) > 1 {
		for _, a := range os.Args[1:] { process(a) }
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { process(sc.Text()) }
}
