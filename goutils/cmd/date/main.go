// date - Display or format the current date and time
// Usage: date [+format] [-d datestring] [-u]
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	utc        = flag.Bool("u", false, "Print UTC time")
	dateString = flag.String("d", "", "Parse and display a date string")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: date [+format] [-d datestring] [-u]")
		fmt.Fprintln(os.Stderr, "Format specifiers: %Y %m %d %H %M %S %A %B %Z %s %R %T %D %F %e %j %u %w %W %p %I")
		flag.PrintDefaults()
	}
	flag.Parse()

	var t time.Time
	if *dateString != "" {
		var err error
		layouts := []string{
			"2006-01-02 15:04:05",
			"2006-01-02",
			"01/02/2006",
			"Jan 2 2006",
			"January 2 2006",
			"Mon Jan 2 15:04:05 2006",
			time.RFC3339,
			time.RFC1123,
			time.RFC822,
		}
		for _, layout := range layouts {
			t, err = time.Parse(layout, *dateString)
			if err == nil {
				break
			}
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "date: cannot parse date:", *dateString)
			os.Exit(1)
		}
	} else {
		t = time.Now()
	}

	if *utc {
		t = t.UTC()
	}

	format := "+%a %b %e %H:%M:%S %Z %Y"
	if flag.NArg() > 0 {
		format = flag.Arg(0)
	}

	if strings.HasPrefix(format, "+") {
		format = format[1:]
	}

	fmt.Println(formatDate(t, format))
}

func formatDate(t time.Time, format string) string {
	weekdays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	months := []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	replacements := map[string]string{
		"%Y": fmt.Sprintf("%04d", t.Year()),
		"%y": fmt.Sprintf("%02d", t.Year()%100),
		"%m": fmt.Sprintf("%02d", t.Month()),
		"%d": fmt.Sprintf("%02d", t.Day()),
		"%e": fmt.Sprintf("%2d", t.Day()),
		"%H": fmt.Sprintf("%02d", t.Hour()),
		"%I": fmt.Sprintf("%02d", (t.Hour()%12 + 12) % 12),
		"%M": fmt.Sprintf("%02d", t.Minute()),
		"%S": fmt.Sprintf("%02d", t.Second()),
		"%A": weekdays[t.Weekday()],
		"%a": weekdays[t.Weekday()][:3],
		"%B": months[t.Month()-1],
		"%b": months[t.Month()-1][:3],
		"%Z": t.Format("MST"),
		"%z": t.Format("-0700"),
		"%s": fmt.Sprintf("%d", t.Unix()),
		"%n": "\n",
		"%t": "\t",
		"%%": "%",
		"%R": fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute()),
		"%T": fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second()),
		"%D": fmt.Sprintf("%02d/%02d/%02d", t.Month(), t.Day(), t.Year()%100),
		"%F": fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day()),
		"%j": fmt.Sprintf("%03d", t.YearDay()),
		"%u": fmt.Sprintf("%d", int(t.Weekday())+1),
		"%w": fmt.Sprintf("%d", int(t.Weekday())),
		"%p": func() string {
			if t.Hour() < 12 {
				return "AM"
			}
			return "PM"
		}(),
	}

	result := format
	for k, v := range replacements {
		result = strings.ReplaceAll(result, k, v)
	}
	return result
}
