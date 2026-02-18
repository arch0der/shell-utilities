// watcher - Watch a directory for file changes
// Usage: watcher [-r] <dir>
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var recursive = flag.Bool("r", false, "Watch subdirectories recursively")

type fileState struct {
	size    int64
	modTime time.Time
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: watcher [-r] <dir>")
		flag.PrintDefaults()
	}
	flag.Parse()

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	fmt.Printf("Watching %s for changes (Ctrl+C to stop)...\n", dir)

	prev := snapshot(dir)
	for {
		time.Sleep(500 * time.Millisecond)
		curr := snapshot(dir)

		for path, cs := range curr {
			ps, existed := prev[path]
			if !existed {
				fmt.Printf("CREATED  %s\n", path)
			} else if cs.modTime != ps.modTime || cs.size != ps.size {
				fmt.Printf("MODIFIED %s\n", path)
			}
		}
		for path := range prev {
			if _, ok := curr[path]; !ok {
				fmt.Printf("DELETED  %s\n", path)
			}
		}
		prev = curr
	}
}

func snapshot(dir string) map[string]fileState {
	state := map[string]fileState{}
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !*recursive && filepath.Dir(path) != dir && path != dir {
			if info.IsDir() {
				return filepath.SkipDir
			}
		}
		if !info.IsDir() {
			state[path] = fileState{size: info.Size(), modTime: info.ModTime()}
		}
		return nil
	}
	filepath.Walk(dir, walk)
	return state
}
