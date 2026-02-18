// cal - Print a calendar
// Usage: cal [month] [year]
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	if len(os.Args) >= 3 {
		month, _ = strconv.Atoi(os.Args[1])
		year, _ = strconv.Atoi(os.Args[2])
	} else if len(os.Args) == 2 {
		year, _ = strconv.Atoi(os.Args[1])
		month = 0 // print whole year
	}

	if month == 0 {
		printYear(year)
	} else {
		printMonth(year, time.Month(month), now)
	}
}

func printMonth(year int, month time.Month, now time.Time) {
	fmt.Printf("    %s %d\n", month.String(), year)
	fmt.Println("Su Mo Tu We Th Fr Sa")

	first := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	weekday := int(first.Weekday())

	// Print leading spaces
	fmt.Print(spaces(weekday * 3))

	day := 1
	col := weekday
	for {
		d := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		if d.Month() != month {
			break
		}
		isToday := d.Year() == now.Year() && d.Month() == now.Month() && d.Day() == now.Day()
		if isToday {
			fmt.Printf("[%2d]", day)
		} else {
			fmt.Printf("%2d ", day)
		}
		col++
		if col%7 == 0 {
			fmt.Println()
		}
		day++
	}
	if col%7 != 0 {
		fmt.Println()
	}
}

func printYear(year int) {
	now := time.Now()
	fmt.Printf("         %d\n\n", year)
	for m := time.January; m <= time.December; m++ {
		printMonth(year, m, now)
		fmt.Println()
	}
}

func spaces(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}
