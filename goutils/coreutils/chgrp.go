package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func init() { register("chgrp", runChgrp) }

func runChgrp() {
	args := os.Args[1:]
	recursive := false
	files := []string{}
	group := ""
	for _, a := range args {
		if a == "-R" || a == "-r" || a == "--recursive" {
			recursive = true
		} else if !strings.HasPrefix(a, "-") {
			if group == "" {
				group = a
			} else {
				files = append(files, a)
			}
		}
	}
	if group == "" || len(files) == 0 {
		fmt.Fprintln(os.Stderr, "chgrp: missing operand")
		os.Exit(1)
	}
	var gid int
	if id, err := strconv.Atoi(group); err == nil {
		gid = id
	} else {
		g, err := user.LookupGroup(group)
		if err != nil {
			fmt.Fprintf(os.Stderr, "chgrp: invalid group: %s\n", group)
			os.Exit(1)
		}
		gid, _ = strconv.Atoi(g.Gid)
	}
	var doChgrp func(path string)
	doChgrp = func(path string) {
		info, err := os.Lstat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "chgrp: %s: %v\n", path, err)
			return
		}
		if err := os.Lchown(path, -1, gid); err != nil {
			fmt.Fprintf(os.Stderr, "chgrp: %s: %v\n", path, err)
		}
		if recursive && info.IsDir() {
			entries, _ := os.ReadDir(path)
			for _, e := range entries {
				doChgrp(filepath.Join(path, e.Name()))
			}
		}
	}
	for _, f := range files {
		doChgrp(f)
	}
}
