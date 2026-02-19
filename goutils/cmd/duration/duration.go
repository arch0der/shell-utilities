// duration - parse and convert durations between units, humanize seconds
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type unit struct{ name string; secs float64 }
var units = []unit{
	{"year", 365.25 * 24 * 3600}, {"week", 7 * 24 * 3600},
	{"day", 24 * 3600}, {"hour", 3600}, {"minute", 60}, {"second", 1},
	{"millisecond", 0.001}, {"microsecond", 0.000001}, {"nanosecond", 0.000000001},
}

var abbrev = map[string]float64{
	"ns": 1e-9, "us": 1e-6, "µs": 1e-6, "ms": 0.001,
	"s": 1, "sec": 1, "secs": 1, "second": 1, "seconds": 1,
	"m": 60, "min": 60, "minute": 60, "minutes": 60,
	"h": 3600, "hr": 3600, "hour": 3600, "hours": 3600,
	"d": 86400, "day": 86400, "days": 86400,
	"w": 604800, "week": 604800, "weeks": 604800,
	"y": 31557600, "year": 31557600, "years": 31557600,
}

func parse(s string) (float64, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	// try pure number (seconds)
	if f, err := strconv.ParseFloat(s, 64); err == nil { return f, nil }
	// split number and unit
	i := 0
	for i < len(s) && (s[i] >= '0' && s[i] <= '9' || s[i] == '.' || s[i] == '-') { i++ }
	numStr, unitStr := strings.TrimSpace(s[:i]), strings.TrimSpace(s[i:])
	f, err := strconv.ParseFloat(numStr, 64)
	if err != nil { return 0, fmt.Errorf("cannot parse %q", s) }
	mult, ok := abbrev[unitStr]
	if !ok { return 0, fmt.Errorf("unknown unit %q", unitStr) }
	return f * mult, nil
}

func humanize(secs float64) string {
	if secs < 0 { return "-" + humanize(-secs) }
	var parts []string
	remaining := secs
	for _, u := range units {
		if remaining >= u.secs {
			count := int(remaining / u.secs)
			remaining -= float64(count) * u.secs
			name := u.name
			if count != 1 { name += "s" }
			parts = append(parts, fmt.Sprintf("%d %s", count, name))
			if len(parts) == 3 { break }
		}
	}
	if len(parts) == 0 { return fmt.Sprintf("%.3f seconds", secs) }
	return strings.Join(parts, ", ")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: duration <value[unit]> [target_unit]")
		fmt.Fprintln(os.Stderr, "  examples: duration 3600 hours  |  duration '2 days' seconds  |  duration 90")
		os.Exit(1)
	}
	secs, err := parse(os.Args[1])
	if err != nil { fmt.Fprintln(os.Stderr, "duration:", err); os.Exit(1) }

	fmt.Printf("Human    : %s\n", humanize(secs))
	if len(os.Args) > 2 {
		to := strings.ToLower(strings.TrimSpace(os.Args[2]))
		mult, ok := abbrev[to]
		if !ok { fmt.Fprintf(os.Stderr, "duration: unknown unit %q\n", to); os.Exit(1) }
		fmt.Printf("In %-10s: %g\n", to, secs/mult)
		return
	}
	for _, u := range units {
		fmt.Printf("  %-15s: %g\n", u.name+"s", secs/u.secs)
	}
}
