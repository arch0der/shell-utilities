// daysbetween - calculate the number of days (and other units) between two dates
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var formats = []string{
	"2006-01-02", "01/02/2006", "02-01-2006", "January 2, 2006",
	"Jan 2, 2006", "2 January 2006", "2006-01-02T15:04:05", time.RFC3339,
}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "today" || s == "now" { return time.Now(), nil }
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil { return t, nil }
	}
	return time.Time{}, fmt.Errorf("cannot parse date %q", s)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: daysbetween <date1> <date2>")
		fmt.Fprintln(os.Stderr, "  dates: YYYY-MM-DD, today, now, MM/DD/YYYY, etc.")
		os.Exit(1)
	}
	a, err := parseDate(os.Args[1]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	b, err := parseDate(os.Args[2]); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

	if b.Before(a) { a, b = b, a }
	diff := b.Sub(a)
	days := int(diff.Hours() / 24)
	weeks := days / 7
	months := 0; tmp := a
	for tmp.Before(b) { tmp = tmp.AddDate(0, 1, 0); months++ }
	if tmp.After(b) { months-- }
	years := months / 12

	fmt.Printf("From  : %s\n", a.Format("Monday, January 2, 2006"))
	fmt.Printf("To    : %s\n", b.Format("Monday, January 2, 2006"))
	fmt.Printf("Days  : %d\n", days)
	fmt.Printf("Weeks : %d weeks + %d days\n", weeks, days%7)
	fmt.Printf("Months: %d months + %d days\n", months, days-months*30)
	fmt.Printf("Years : %d years, %d months, %d days\n", years, months%12, days%30)
	fmt.Printf("Hours : %d\n", int(diff.Hours()))
	fmt.Printf("Weekdays: %d\n", countWeekdays(a, b))
}

func countWeekdays(start, end time.Time) int {
	count := 0
	cur := start
	for cur.Before(end) {
		if cur.Weekday() != time.Saturday && cur.Weekday() != time.Sunday { count++ }
		cur = cur.AddDate(0, 0, 1)
	}
	return count
}
