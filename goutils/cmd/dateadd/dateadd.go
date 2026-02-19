// dateadd - add or subtract durations from dates
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const helpMsg = `usage: dateadd <date> <+/-amount> <unit> [format]
  date   : YYYY-MM-DD or "now"
  unit   : seconds|minutes|hours|days|weeks|months|years
  format : output format (default: 2006-01-02)
examples:
  dateadd now +3 days
  dateadd 2024-01-15 -2 months
  dateadd now +1 years "Jan 2, 2006"`

func main() {
	if len(os.Args) < 4 { fmt.Fprintln(os.Stderr, helpMsg); os.Exit(1) }

	dateStr := os.Args[1]
	amtStr := os.Args[2]
	unit := strings.ToLower(os.Args[3])
	outFmt := "2006-01-02"
	if len(os.Args) > 4 { outFmt = os.Args[4] }

	var base time.Time
	var err error
	if dateStr == "now" || dateStr == "today" {
		base = time.Now()
	} else {
		for _, f := range []string{"2006-01-02", "01/02/2006", "2006-01-02 15:04:05", time.RFC3339} {
			base, err = time.Parse(f, dateStr)
			if err == nil { break }
		}
		if err != nil { fmt.Fprintf(os.Stderr, "dateadd: cannot parse date %q\n", dateStr); os.Exit(1) }
	}

	amt, err := strconv.Atoi(amtStr)
	if err != nil { fmt.Fprintf(os.Stderr, "dateadd: invalid amount %q\n", amtStr); os.Exit(1) }

	var result time.Time
	switch unit {
	case "second", "seconds", "s": result = base.Add(time.Duration(amt) * time.Second)
	case "minute", "minutes", "m": result = base.Add(time.Duration(amt) * time.Minute)
	case "hour", "hours", "h": result = base.Add(time.Duration(amt) * time.Hour)
	case "day", "days", "d": result = base.AddDate(0, 0, amt)
	case "week", "weeks", "w": result = base.AddDate(0, 0, amt*7)
	case "month", "months": result = base.AddDate(0, amt, 0)
	case "year", "years", "y": result = base.AddDate(amt, 0, 0)
	default: fmt.Fprintf(os.Stderr, "dateadd: unknown unit %q\n", unit); os.Exit(1)
	}

	fmt.Println(result.Format(outFmt))
}
