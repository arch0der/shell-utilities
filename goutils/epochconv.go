// epochconv - Convert between Unix timestamps and human-readable dates.
//
// Usage:
//
//	epochconv [OPTIONS] [VALUE...]
//	echo 1700000000 | epochconv
//
// Options:
//
//	-f FORMAT  Output format (Go layout, default: "2006-01-02 15:04:05 UTC")
//	-z ZONE    Timezone (default: UTC, e.g. America/New_York)
//	-ms        Treat input as milliseconds
//	-us        Treat input as microseconds
//	-r         Reverse: parse date string to epoch
//	-n         Print current epoch timestamp
//
// Examples:
//
//	epochconv 1700000000                   # 2023-11-14 22:13:20 UTC
//	epochconv -z America/New_York 1700000000
//	epochconv -r "2023-11-14 22:13:20"     # 1700000000
//	epochconv -n                            # current epoch
//	epochconv -ms 1700000000000            # milliseconds input
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	format   = flag.String("f", "2006-01-02 15:04:05 UTC", "output format")
	zone     = flag.String("z", "UTC", "timezone")
	millis   = flag.Bool("ms", false, "millisecond input")
	micros   = flag.Bool("us", false, "microsecond input")
	reverse  = flag.Bool("r", false, "date string to epoch")
	now      = flag.Bool("n", false, "print current epoch")
)

func loadZone() *time.Location {
	loc, err := time.LoadLocation(*zone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "epochconv: unknown timezone %q\n", *zone)
		os.Exit(1)
	}
	return loc
}

func epochToTime(s string, loc *time.Location) (time.Time, error) {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	if *micros {
		return time.Unix(n/1e6, (n%1e6)*1000).In(loc), nil
	}
	if *millis {
		return time.Unix(n/1000, (n%1000)*1e6).In(loc), nil
	}
	return time.Unix(n, 0).In(loc), nil
}

var parseFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02",
	"01/02/2006",
	"Jan 2, 2006",
	"January 2, 2006",
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
}

func parseDate(s string, loc *time.Location) (time.Time, error) {
	s = strings.TrimSpace(s)
	for _, f := range parseFormats {
		t, err := time.ParseInLocation(f, s, loc)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse %q", s)
}

func main() {
	flag.Parse()
	loc := loadZone()
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if *now {
		fmt.Fprintf(w, "%d\n", time.Now().Unix())
		return
	}

	process := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		if *reverse {
			t, err := parseDate(s, loc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "epochconv: %v\n", err)
				return
			}
			fmt.Fprintf(w, "%d\n", t.Unix())
		} else {
			t, err := epochToTime(s, loc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "epochconv: not a valid epoch: %q\n", s)
				return
			}
			fmt.Fprintln(w, t.Format(*format))
		}
	}

	args := flag.Args()
	if len(args) > 0 {
		for _, a := range args {
			process(a)
		}
		return
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		process(sc.Text())
	}
}
