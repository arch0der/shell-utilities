package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

func main() {
	args := os.Args[1:]
	humanReadable := false
	files := []string{}
	for _, a := range args {
		switch a {
		case "-h", "--human-readable":
			humanReadable = true
		case "-H", "--si":
			humanReadable = true
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}
	if len(files) == 0 {
		files = []string{"/"}
	}
	fmt.Printf("%-20s %12s %12s %12s %6s %s\n", "Filesystem", "1K-blocks", "Used", "Available", "Use%", "Mounted on")
	for _, f := range files {
		var stat syscall.Statfs_t
		if err := syscall.Statfs(f, &stat); err != nil {
			fmt.Fprintf(os.Stderr, "df: %s: %v\n", f, err)
			continue
		}
		total := int64(stat.Blocks) * int64(stat.Bsize)
		avail := int64(stat.Bavail) * int64(stat.Bsize)
		used := total - int64(stat.Bfree)*int64(stat.Bsize)
		pct := 0
		if total > 0 {
			pct = int(100 * used / total)
		}
		if humanReadable {
			fmt.Printf("%-20s %12s %12s %12s %5d%% %s\n",
				f, humanSize(total), humanSize(used), humanSize(avail), pct, f)
		} else {
			fmt.Printf("%-20s %12d %12d %12d %5d%% %s\n",
				f, total/1024, used/1024, avail/1024, pct, f)
		}
	}
}
