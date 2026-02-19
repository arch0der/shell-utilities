// df - Report disk space usage
// Usage: df [-h] [path...]
package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
)

var human = flag.Bool("h", false, "Human-readable sizes")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: df [-h] [path...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"/"}
	}

	fmt.Printf("%-20s %10s %10s %10s %6s %s\n",
		"Filesystem", "1K-blocks", "Used", "Available", "Use%", "Mounted on")

	for _, path := range paths {
		var stat syscall.Statfs_t
		if err := syscall.Statfs(path, &stat); err != nil {
			fmt.Fprintln(os.Stderr, "df:", err)
			continue
		}

		total := stat.Blocks * uint64(stat.Bsize)
		free := stat.Bfree * uint64(stat.Bsize)
		used := total - free
		avail := stat.Bavail * uint64(stat.Bsize)
		var pct float64
		if total > 0 {
			pct = float64(used) / float64(total) * 100
		}

		if *human {
			fmt.Printf("%-20s %10s %10s %10s %5.0f%% %s\n",
				path, fmtSize(total), fmtSize(used), fmtSize(avail), pct, path)
		} else {
			fmt.Printf("%-20s %10d %10d %10d %5.0f%% %s\n",
				path, total/1024, used/1024, avail/1024, pct, path)
		}
	}
}

func fmtSize(b uint64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1fG", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1fM", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1fK", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%dB", b)
	}
}
