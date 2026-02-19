// cron2human - Translate cron expressions to plain English.
//
// Usage:
//
//	cron2human EXPRESSION
//	echo "0 9 * * 1-5" | cron2human
//
// Options:
//
//	-n N      Show next N execution times (default: 5)
//	-f FORMAT Datetime format (default: "2006-01-02 15:04 MST")
//	-j        JSON output
//
// Supports standard 5-field cron: min hour day month weekday
// Also supports @yearly @monthly @weekly @daily @hourly @reboot
//
// Examples:
//
//	cron2human "0 * * * *"          # every hour
//	cron2human "0 9 * * 1-5"        # weekdays at 9am
//	cron2human "*/15 * * * *"       # every 15 minutes
//	cron2human "@daily"             # every day at midnight
//	cron2human -n 3 "0 0 1 * *"     # next 3 runs
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	nextN  = flag.Int("n", 5, "show next N times")
	format = flag.String("f", "2006-01-02 15:04 MST", "time format")
	asJSON = flag.Bool("j", false, "JSON output")
)

var shortcuts = map[string]string{
	"@yearly":   "0 0 1 1 *",
	"@annually": "0 0 1 1 *",
	"@monthly":  "0 0 1 * *",
	"@weekly":   "0 0 * * 0",
	"@daily":    "0 0 * * *",
	"@midnight": "0 0 * * *",
	"@hourly":   "0 * * * *",
}

var months = []string{"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December"}
var weekdays = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

type Field struct {
	raw   string
	every bool // */N
	list  bool // a,b,c
	rang  bool // a-b
	any   bool // *
	vals  []int
	step  int
}

func parseField(s string, min, max int) Field {
	f := Field{raw: s}
	if s == "*" {
		f.any = true
		return f
	}
	if strings.HasPrefix(s, "*/") {
		n, _ := strconv.Atoi(s[2:])
		f.every = true
		f.step = n
		return f
	}
	if strings.Contains(s, ",") {
		f.list = true
		for _, p := range strings.Split(s, ",") {
			n, _ := strconv.Atoi(p)
			f.vals = append(f.vals, n)
		}
		return f
	}
	if strings.Contains(s, "-") {
		f.rang = true
		parts := strings.SplitN(s, "-", 2)
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		for i := a; i <= b; i++ {
			f.vals = append(f.vals, i)
		}
		return f
	}
	n, _ := strconv.Atoi(s)
	f.vals = []int{n}
	return f
}

func describeField(f Field, unit string, names []string) string {
	if f.any {
		return "every " + unit
	}
	if f.every {
		return fmt.Sprintf("every %d %ss", f.step, unit)
	}
	stringify := func(v int) string {
		if names != nil && v < len(names) {
			return names[v]
		}
		return strconv.Itoa(v)
	}
	parts := make([]string, len(f.vals))
	for i, v := range f.vals {
		parts[i] = stringify(v)
	}
	if len(parts) == 1 {
		return parts[0]
	}
	if f.rang && len(parts) > 2 {
		return parts[0] + " through " + parts[len(parts)-1]
	}
	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
}

func describe(expr string) string {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return "invalid cron expression"
	}

	min := parseField(parts[0], 0, 59)
	hour := parseField(parts[1], 0, 23)
	day := parseField(parts[2], 1, 31)
	month := parseField(parts[3], 1, 12)
	wday := parseField(parts[4], 0, 6)

	var desc strings.Builder

	// Minute
	if min.any {
		desc.WriteString("every minute")
	} else {
		desc.WriteString("at minute " + describeField(min, "minute", nil))
	}

	// Hour
	if !hour.any {
		desc.WriteString(" of hour " + describeField(hour, "hour", nil))
	}

	// Time shorthand
	if len(min.vals) == 1 && len(hour.vals) == 1 {
		h, m := hour.vals[0], min.vals[0]
		ampm := "AM"
		if h >= 12 {
			ampm = "PM"
			if h > 12 {
				h -= 12
			}
		}
		if h == 0 {
			h = 12
		}
		desc.Reset()
		desc.WriteString(fmt.Sprintf("at %d:%02d %s", h, m, ampm))
	}

	// Weekday
	if !wday.any {
		desc.WriteString(", on " + describeField(wday, "weekday", weekdays))
	}

	// Day of month
	if !day.any {
		desc.WriteString(", on day " + describeField(day, "day", nil))
	}

	// Month
	if !month.any {
		// month names are 1-based
		names := append([]string{""}, months...)
		desc.WriteString(", in " + describeField(month, "month", names))
	} else {
		desc.WriteString(", every month")
	}

	return desc.String()
}

func matchesField(f Field, v int) bool {
	if f.any {
		return true
	}
	if f.every {
		return v%f.step == 0
	}
	for _, val := range f.vals {
		if val == v {
			return true
		}
	}
	return false
}

func nextRuns(expr string, n int) []time.Time {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil
	}
	min := parseField(parts[0], 0, 59)
	hour := parseField(parts[1], 0, 23)
	day := parseField(parts[2], 1, 31)
	month := parseField(parts[3], 1, 12)
	wday := parseField(parts[4], 0, 6)

	var results []time.Time
	t := time.Now().Truncate(time.Minute).Add(time.Minute)

	for len(results) < n {
		if matchesField(month, int(t.Month())) &&
			matchesField(day, t.Day()) &&
			matchesField(wday, int(t.Weekday())) &&
			matchesField(hour, t.Hour()) &&
			matchesField(min, t.Minute()) {
			results = append(results, t)
		}
		t = t.Add(time.Minute)
		if t.Year() > time.Now().Year()+5 {
			break
		}
	}
	return results
}

func main() {
	flag.Parse()
	args := flag.Args()

	var expr string
	if len(args) > 0 {
		expr = strings.Join(args, " ")
	} else {
		var b strings.Builder
		fmt.Fscanln(os.Stdin, &b)
		expr = strings.TrimSpace(b.String())
	}

	if sc, ok := shortcuts[expr]; ok {
		expr = sc
	}

	if expr == "@reboot" {
		fmt.Println("at system reboot")
		return
	}

	human := describe(expr)
	runs := nextRuns(expr, *nextN)

	if *asJSON {
		runStrs := make([]string, len(runs))
		for i, r := range runs {
			runStrs[i] = r.Format(*format)
		}
		out := map[string]interface{}{
			"expression":  expr,
			"description": human,
			"next_runs":   runStrs,
		}
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(b))
		return
	}

	fmt.Printf("Expression:  %s\n", expr)
	fmt.Printf("Description: %s\n", human)
	if *nextN > 0 {
		fmt.Printf("Next %d runs:\n", *nextN)
		for _, r := range runs {
			fmt.Printf("  %s\n", r.Format(*format))
		}
	}
}
