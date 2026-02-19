// watch - Execute a program periodically, showing output
// Usage: watch [-n seconds] [-d] <command>
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	interval = flag.Float64("n", 2.0, "Seconds between updates")
	diff     = flag.Bool("d", false, "Highlight differences between updates")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: watch [-n seconds] [-d] <command>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	cmd := strings.Join(flag.Args(), " ")
	prev := ""
	iteration := 0

	for {
		// Clear screen (ANSI escape)
		if iteration > 0 {
			fmt.Print("\033[H\033[2J")
		}

		now := time.Now().Format("Mon Jan  2 15:04:05 2006")
		fmt.Printf("Every %.1fs: %-40s  %s\n\n", *interval, cmd, now)

		out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
		result := string(out)

		if *diff && iteration > 0 {
			// Highlight lines that changed
			prevLines := strings.Split(prev, "\n")
			currLines := strings.Split(result, "\n")
			for i, line := range currLines {
				if i < len(prevLines) && line != prevLines[i] {
					fmt.Printf("\033[7m%s\033[m\n", line) // reverse video
				} else {
					fmt.Println(line)
				}
			}
		} else {
			fmt.Print(result)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "\nCommand exited with: %v\n", err)
		}

		prev = result
		iteration++
		time.Sleep(time.Duration(*interval * float64(time.Second)))
	}
}
