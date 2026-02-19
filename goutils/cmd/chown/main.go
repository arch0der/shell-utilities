// chown - Change file owner and group
// Usage: chown [-R] user[:group] file...
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var recursive = flag.Bool("R", false, "Recursively change ownership")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: chown [-R] user[:group] file...")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	spec := flag.Arg(0)
	uid, gid, err := parseOwner(spec)
	if err != nil {
		fmt.Fprintln(os.Stderr, "chown:", err)
		os.Exit(1)
	}

	exitCode := 0
	for _, path := range flag.Args()[1:] {
		if err := chownPath(path, uid, gid); err != nil {
			fmt.Fprintln(os.Stderr, "chown:", err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func parseOwner(spec string) (uid, gid int, err error) {
	uid, gid = -1, -1
	parts := strings.SplitN(spec, ":", 2)
	userPart := parts[0]
	groupPart := ""
	if len(parts) == 2 {
		groupPart = parts[1]
	}

	if userPart != "" {
		if n, err2 := strconv.Atoi(userPart); err2 == nil {
			uid = n
		} else {
			u, err2 := user.Lookup(userPart)
			if err2 != nil {
				return 0, 0, fmt.Errorf("invalid user: %s", userPart)
			}
			uid, _ = strconv.Atoi(u.Uid)
			// Default group to user's group if not specified
			if groupPart == "" {
				gid, _ = strconv.Atoi(u.Gid)
			}
		}
	}

	if groupPart != "" {
		if n, err2 := strconv.Atoi(groupPart); err2 == nil {
			gid = n
		} else {
			g, err2 := user.LookupGroup(groupPart)
			if err2 != nil {
				return 0, 0, fmt.Errorf("invalid group: %s", groupPart)
			}
			gid, _ = strconv.Atoi(g.Gid)
		}
	}
	return uid, gid, nil
}

func chownPath(path string, uid, gid int) error {
	if *recursive {
		return filepath.Walk(path, func(p string, _ os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			return syscall.Lchown(p, uid, gid)
		})
	}
	return syscall.Lchown(path, uid, gid)
}
