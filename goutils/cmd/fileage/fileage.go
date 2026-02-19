// fileage - show how old files are in human-readable form
package main

import (
	"fmt"
	"os"
	"sort"
	"time"
)

func humanAge(d time.Duration) string {
	switch {
	case d < time.Minute: return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour: return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour: return fmt.Sprintf("%.1fh ago", d.Hours())
	case d < 7*24*time.Hour: return fmt.Sprintf("%.1f days ago", d.Hours()/24)
	case d < 365*24*time.Hour: return fmt.Sprintf("%.1f weeks ago", d.Hours()/(24*7))
	default: return fmt.Sprintf("%.1f years ago", d.Hours()/(24*365.25))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: fileage <file> [file...]")
		os.Exit(1)
	}

	type entry struct {
		name string
		mod  time.Time
		age  time.Duration
	}
	now := time.Now()
	var entries []entry
	for _, path := range os.Args[1:] {
		info, err := os.Stat(path)
		if err != nil { fmt.Fprintln(os.Stderr, err); continue }
		mod := info.ModTime()
		entries = append(entries, entry{path, mod, now.Sub(mod)})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].age < entries[j].age })
	fmt.Printf("%-40s  %-28s  %s\n", "File", "Modified", "Age")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────")
	for _, e := range entries {
		fmt.Printf("%-40s  %-28s  %s\n", e.name, e.mod.Format("2006-01-02 15:04:05 MST"), humanAge(e.age))
	}
}
