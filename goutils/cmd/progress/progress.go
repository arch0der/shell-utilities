// progress - show a progress bar while piping data or counting lines
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func humanSize(b int64) string {
	switch {
	case b < 1024: return fmt.Sprintf("%dB", b)
	case b < 1024*1024: return fmt.Sprintf("%.1fKB", float64(b)/1024)
	case b < 1024*1024*1024: return fmt.Sprintf("%.1fMB", float64(b)/(1024*1024))
	default: return fmt.Sprintf("%.2fGB", float64(b)/(1024*1024*1024))
	}
}

func drawBar(done, total int64, elapsed time.Duration, width int) string {
	pct := 0.0
	if total > 0 { pct = float64(done) / float64(total) }
	filled := int(pct * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	speed := float64(done) / elapsed.Seconds()
	eta := ""
	if total > 0 && speed > 0 {
		secs := float64(total-done) / speed
		eta = fmt.Sprintf("ETA %ds", int(secs))
	}
	if total > 0 {
		return fmt.Sprintf("\r[%s] %.1f%%  %s/%s  %.0f/s  %s",
			bar, pct*100, humanSize(done), humanSize(total), speed, eta)
	}
	return fmt.Sprintf("\r[%s]  %s  %.0f/s  %s",
		bar[:filled+1]+"...", humanSize(done), speed, elapsed.Round(time.Second))
}

func main() {
	total := int64(-1)
	lineMode := false
	width := 40
	interval := 100 * time.Millisecond

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-s", "--size": i++; s, _ := strconv.ParseInt(args[i], 10, 64); total = s
		case "-l", "--lines": lineMode = true
		case "-n", "--total": i++; total, _ = strconv.ParseInt(args[i], 10, 64)
		case "-w": i++; width, _ = strconv.Atoi(args[i])
		}
	}

	start := time.Now()
	lastDraw := time.Now()
	var count int64

	if lineMode {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			fmt.Fprintln(os.Stdout, sc.Text())
			count++
			if time.Since(lastDraw) > interval {
				lastDraw = time.Now()
				elapsed := time.Since(start)
				bar := drawBar(count, total, elapsed, width)
				fmt.Fprint(os.Stderr, bar)
			}
		}
	} else {
		buf := make([]byte, 32*1024)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				os.Stdout.Write(buf[:n])
				count += int64(n)
				if time.Since(lastDraw) > interval {
					lastDraw = time.Now()
					fmt.Fprint(os.Stderr, drawBar(count, total, time.Since(start), width))
				}
			}
			if err == io.EOF { break }
		}
	}

	elapsed := time.Since(start)
	fmt.Fprintf(os.Stderr, "\r\033[K") // clear line
	label := humanSize(count); if lineMode { label = fmt.Sprintf("%d lines", count) }
	speed := float64(count) / elapsed.Seconds()
	fmt.Fprintf(os.Stderr, "Done: %s in %s (%.0f/s)\n", label, elapsed.Round(time.Millisecond), speed)
}
