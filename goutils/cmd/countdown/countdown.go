// countdown - count down from N seconds, showing a live timer
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: countdown <seconds> [label]")
		os.Exit(1)
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n < 0 {
		fmt.Fprintln(os.Stderr, "countdown: seconds must be a non-negative integer")
		os.Exit(1)
	}
	label := "Time remaining"
	if len(os.Args) > 2 { label = os.Args[2] }

	for i := n; i >= 0; i-- {
		h := i / 3600
		m := (i % 3600) / 60
		s := i % 60
		fmt.Printf("\r%s: %02d:%02d:%02d  ", label, h, m, s)
		if i == 0 { break }
		time.Sleep(time.Second)
	}
	fmt.Println("\n\a🔔 Done!")
}
