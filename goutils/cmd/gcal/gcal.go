// gcal - print a graphical calendar for a month or year
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var monthNames = []string{"", "January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December"}

func printMonth(year, month int, highlight int) {
	first := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.Local).Day()
	startDay := int(first.Weekday()) // 0=Sun

	title := fmt.Sprintf("%s %d", monthNames[month], year)
	pad := (20 - len(title)) / 2
	fmt.Printf("%s%s\n", strings.Repeat(" ", pad), title)
	fmt.Println("Su Mo Tu We Th Fr Sa")

	// Leading spaces
	fmt.Print(strings.Repeat("   ", startDay))
	col := startDay
	for d := 1; d <= daysInMonth; d++ {
		if d == highlight {
			fmt.Printf("\033[7m%2d\033[0m", d)
		} else {
			fmt.Printf("%2d", d)
		}
		col++
		if col == 7 { fmt.Println(); col = 0 } else { fmt.Print(" ") }
	}
	if col != 0 { fmt.Println() }
	fmt.Println()
}

func printYear(year int) {
	fmt.Printf("%s%d\n\n", strings.Repeat(" ", 30), year)
	for row := 0; row < 4; row++ {
		// Print 3 months side by side
		months := make([][]string, 3)
		for col := 0; col < 3; col++ {
			m := row*3 + col + 1
			months[col] = renderMonth(year, m, -1)
		}
		maxLines := 0
		for _, m := range months { if len(m) > maxLines { maxLines = len(m) } }
		for line := 0; line < maxLines; line++ {
			for col, m := range months {
				s := ""
				if line < len(m) { s = m[line] }
				fmt.Printf("%-22s", s)
				if col < 2 { fmt.Print("  ") }
			}
			fmt.Println()
		}
	}
}

func renderMonth(year, month, highlight int) []string {
	var lines []string
	first := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.Local).Day()
	startDay := int(first.Weekday())
	title := fmt.Sprintf("%s %d", monthNames[month][:3], year)
	pad := (20 - len(title)) / 2
	lines = append(lines, strings.Repeat(" ", pad)+title)
	lines = append(lines, "Su Mo Tu We Th Fr Sa")
	row := strings.Repeat("   ", startDay)
	col := startDay
	for d := 1; d <= daysInMonth; d++ {
		if d == highlight { row += fmt.Sprintf("\033[7m%2d\033[0m", d) } else { row += fmt.Sprintf("%2d", d) }
		col++
		if col == 7 { lines = append(lines, row); row = ""; col = 0 } else { row += " " }
	}
	if col != 0 { lines = append(lines, row) }
	return lines
}

func main() {
	now := time.Now()
	year, month := now.Year(), int(now.Month())
	today := now.Day()
	yearMode := false

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-y", "--year": yearMode = true
		default:
			if n, err := strconv.Atoi(args[i]); err == nil {
				if n > 12 { year = n } else { month = n }
			}
		}
	}

	if yearMode { printYear(year); return }
	printMonth(year, month, func() int {
		if year == now.Year() && month == int(now.Month()) { return today }
		return -1
	}())
}
