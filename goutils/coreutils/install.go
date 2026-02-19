package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func init() { register("install", runInstall) }

func runInstall() {
	args := os.Args[1:]
	mode := os.FileMode(0755)
	ownerStr := ""
	groupStr := ""
	dirMode := false
	backupSuffix := ""
	verbose := false
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-d" || a == "--directory":
			dirMode = true
		case a == "-v" || a == "--verbose":
			verbose = true
		case a == "-m" && i+1 < len(args):
			i++
			m, err := strconv.ParseUint(args[i], 8, 32)
			if err == nil {
				mode = os.FileMode(m)
			}
		case strings.HasPrefix(a, "-m"):
			m, _ := strconv.ParseUint(a[2:], 8, 32)
			mode = os.FileMode(m)
		case a == "-o" && i+1 < len(args):
			i++
			ownerStr = args[i]
		case a == "-g" && i+1 < len(args):
			i++
			groupStr = args[i]
		case a == "-b":
			backupSuffix = "~"
		case strings.HasPrefix(a, "--suffix="):
			backupSuffix = a[9:]
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}

	uid, gid := -1, -1
	if ownerStr != "" {
		if id, err := strconv.Atoi(ownerStr); err == nil {
			uid = id
		} else if u, err := user.Lookup(ownerStr); err == nil {
			uid, _ = strconv.Atoi(u.Uid)
		}
	}
	if groupStr != "" {
		if id, err := strconv.Atoi(groupStr); err == nil {
			gid = id
		} else if g, err := user.LookupGroup(groupStr); err == nil {
			gid, _ = strconv.Atoi(g.Gid)
		}
	}

	if dirMode {
		for _, d := range files {
			if err := os.MkdirAll(d, mode); err != nil {
				fmt.Fprintf(os.Stderr, "install: %s: %v\n", d, err)
			} else {
				if verbose {
					fmt.Printf("install: creating directory '%s'\n", d)
				}
				if uid >= 0 || gid >= 0 {
					os.Chown(d, uid, gid)
				}
			}
		}
		return
	}

	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "install: missing destination")
		os.Exit(1)
	}
	dest := files[len(files)-1]
	srcs := files[:len(files)-1]

	for _, src := range srcs {
		dst := dest
		if info, err := os.Stat(dest); err == nil && info.IsDir() {
			dst = filepath.Join(dest, filepath.Base(src))
		}
		if backupSuffix != "" {
			if _, err := os.Stat(dst); err == nil {
				os.Rename(dst, dst+backupSuffix)
			}
		}
		in, err := os.Open(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "install: %v\n", err)
			continue
		}
		out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "install: %v\n", err)
			in.Close()
			continue
		}
		io.Copy(out, in)
		in.Close()
		out.Close()
		if uid >= 0 || gid >= 0 {
			os.Chown(dst, uid, gid)
		}
		if verbose {
			fmt.Printf("'%s' -> '%s'\n", src, dst)
		}
	}
}
