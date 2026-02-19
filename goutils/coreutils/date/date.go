package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	format := ""
	utc := false
	setDate := ""

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-u" || a == "--utc" || a == "--universal":
			utc = true
		case a == "-d" && i+1 < len(args):
			i++
			setDate = args[i]
		case strings.HasPrefix(a, "-d"):
			setDate = a[2:]
		case strings.HasPrefix(a, "+"):
			format = a[1:]
		}
	}

	var t time.Time
	if setDate != "" {
		var err error
		layouts := []string{
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
			"Mon Jan 2 15:04:05 MST 2006",
			"January 2, 2006",
			"2006-01-02",
		}
		for _, l := range layouts {
			t, err = time.Parse(l, setDate)
			if err == nil {
				break
			}
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "date: invalid date '%s'\n", setDate)
			os.Exit(1)
		}
	} else {
		t = time.Now()
	}

	if utc {
		t = t.UTC()
	}

	if format == "" {
		// Default format: Mon Jan  2 15:04:05 MST 2006
		fmt.Println(t.Format("Mon Jan  2 15:04:05 MST 2006"))
		return
	}

	// Convert strftime format to Go format
	result := strftimeFormat(format, t)
	fmt.Println(result)
}

func strftimeFormat(format string, t time.Time) string {
	var sb strings.Builder
	for i := 0; i < len(format); i++ {
		if format[i] != '%' || i+1 >= len(format) {
			sb.WriteByte(format[i])
			continue
		}
		i++
		switch format[i] {
		case 'Y':
			sb.WriteString(t.Format("2006"))
		case 'y':
			sb.WriteString(t.Format("06"))
		case 'm':
			sb.WriteString(t.Format("01"))
		case 'd':
			sb.WriteString(t.Format("02"))
		case 'H':
			sb.WriteString(t.Format("15"))
		case 'M':
			sb.WriteString(t.Format("04"))
		case 'S':
			sb.WriteString(t.Format("05"))
		case 'A':
			sb.WriteString(t.Format("Monday"))
		case 'a':
			sb.WriteString(t.Format("Mon"))
		case 'B':
			sb.WriteString(t.Format("January"))
		case 'b', 'h':
			sb.WriteString(t.Format("Jan"))
		case 'I':
			sb.WriteString(t.Format("03"))
		case 'p':
			sb.WriteString(t.Format("PM"))
		case 'P':
			sb.WriteString(strings.ToLower(t.Format("PM")))
		case 'Z':
			sb.WriteString(t.Format("MST"))
		case 'z':
			sb.WriteString(t.Format("-0700"))
		case 'j':
			sb.WriteString(fmt.Sprintf("%03d", t.YearDay()))
		case 'u':
			d := int(t.Weekday())
			if d == 0 {
				d = 7
			}
			sb.WriteString(fmt.Sprintf("%d", d))
		case 'w':
			sb.WriteString(fmt.Sprintf("%d", int(t.Weekday())))
		case 'e':
			sb.WriteString(fmt.Sprintf("%2d", t.Day()))
		case 'n':
			sb.WriteByte('\n')
		case 't':
			sb.WriteByte('\t')
		case '%':
			sb.WriteByte('%')
		case 's':
			sb.WriteString(fmt.Sprintf("%d", t.Unix()))
		default:
			sb.WriteByte('%')
			sb.WriteByte(format[i])
		}
	}
	return sb.String()
}
