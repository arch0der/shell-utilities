package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	recursive := false
	ownerStr := ""
	files := []string{}
	for _, a := range args {
		if a == "-R" || a == "--recursive" {
			recursive = true
		} else if !strings.HasPrefix(a, "-") {
			if ownerStr == "" {
				ownerStr = a
			} else {
				files = append(files, a)
			}
		}
	}
	if ownerStr == "" || len(files) == 0 {
		fmt.Fprintln(os.Stderr, "chown: missing operand")
		os.Exit(1)
	}
	parts := strings.SplitN(ownerStr, ":", 2)
	uid, gid := -1, -1
	if parts[0] != "" {
		if id, err := strconv.Atoi(parts[0]); err == nil {
			uid = id
		} else {
			u, err := user.Lookup(parts[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "chown: invalid user: %s\n", parts[0])
				os.Exit(1)
			}
			uid, _ = strconv.Atoi(u.Uid)
		}
	}
	if len(parts) == 2 && parts[1] != "" {
		if id, err := strconv.Atoi(parts[1]); err == nil {
			gid = id
		} else {
			g, err := user.LookupGroup(parts[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "chown: invalid group: %s\n", parts[1])
				os.Exit(1)
			}
			gid, _ = strconv.Atoi(g.Gid)
		}
	}
	var doChown func(path string)
	doChown = func(path string) {
		if err := os.Lchown(path, uid, gid); err != nil {
			fmt.Fprintf(os.Stderr, "chown: %s: %v\n", path, err)
		}
		if recursive {
			info, _ := os.Lstat(path)
			if info != nil && info.IsDir() {
				entries, _ := os.ReadDir(path)
				for _, e := range entries {
					doChown(filepath.Join(path, e.Name()))
				}
			}
		}
	}
	for _, f := range files {
		doChown(f)
	}
}
