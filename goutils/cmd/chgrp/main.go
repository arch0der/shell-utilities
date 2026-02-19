// chgrp - Change group ownership of files
// Usage: chgrp [-R] group file...
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

var recursive = flag.Bool("R", false, "Recursively change group")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: chgrp [-R] group file...")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	groupSpec := flag.Arg(0)
	gid, err := resolveGroup(groupSpec)
	if err != nil {
		fmt.Fprintln(os.Stderr, "chgrp:", err)
		os.Exit(1)
	}

	exitCode := 0
	for _, path := range flag.Args()[1:] {
		if err := chgrpPath(path, gid); err != nil {
			fmt.Fprintln(os.Stderr, "chgrp:", err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func resolveGroup(spec string) (int, error) {
	if n, err := strconv.Atoi(spec); err == nil {
		return n, nil
	}
	g, err := user.LookupGroup(spec)
	if err != nil {
		return 0, fmt.Errorf("invalid group: %s", spec)
	}
	gid, _ := strconv.Atoi(g.Gid)
	return gid, nil
}

func chgrpPath(path string, gid int) error {
	if *recursive {
		return filepath.Walk(path, func(p string, _ os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			info, err := os.Lstat(p)
			if err != nil {
				return err
			}
			stat := info.Sys().(*syscall.Stat_t)
			return syscall.Lchown(p, int(stat.Uid), gid)
		})
	}
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	stat := info.Sys().(*syscall.Stat_t)
	return syscall.Lchown(path, int(stat.Uid), gid)
}
