// humanize - convert raw numbers to human-readable sizes, counts, and durations
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func humanBytes(n float64) string {
	units := []string{"B","KB","MB","GB","TB","PB"}
	i := 0
	for n >= 1024 && i < len(units)-1 { n /= 1024; i++ }
	if i == 0 { return fmt.Sprintf("%.0f %s", n, units[i]) }
	return fmt.Sprintf("%.2f %s", n, units[i])
}

func humanCount(n float64) string {
	switch {
	case n >= 1e12: return fmt.Sprintf("%.2fT", n/1e12)
	case n >= 1e9: return fmt.Sprintf("%.2fB", n/1e9)
	case n >= 1e6: return fmt.Sprintf("%.2fM", n/1e6)
	case n >= 1e3: return fmt.Sprintf("%.2fK", n/1e3)
	default: return fmt.Sprintf("%.0f", n)
	}
}

func humanDuration(secs float64) string {
	if secs < 60 { return fmt.Sprintf("%.1fs", secs) }
	if secs < 3600 { return fmt.Sprintf("%.1fm", secs/60) }
	if secs < 86400 { return fmt.Sprintf("%.1fh", secs/3600) }
	return fmt.Sprintf("%.1fd", secs/86400)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: humanize [bytes|count|duration] [numbers...  |  stdin]")
	os.Exit(1)
}

func main() {
	mode := "bytes"
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "bytes","count","duration": mode = args[0]; args = args[1:]
		}
	}

	fn := map[string]func(float64)string{
		"bytes": humanBytes, "count": humanCount, "duration": humanDuration,
	}[mode]

	process := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" { return }
		f, err := strconv.ParseFloat(s, 64)
		if err != nil { fmt.Fprintf(os.Stderr, "humanize: %q not a number\n", s); return }
		fmt.Printf("%-20g → %s\n", f, fn(f))
	}

	if len(args) > 0 { for _, a := range args { process(a) }; return }
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { process(sc.Text()) }
}
