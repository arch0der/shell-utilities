package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	args := os.Args[1:]
	long := false
	all := false
	almostAll := false
	humanReadable := false
	reverse := false
	sortByTime := false
	sortBySize := false
	noSort := false
	recursive := false
	inode := false
	classify := false
	colorize := false
	onePerLine := false
	dirs := []string{}

	for _, a := range args {
		if strings.HasPrefix(a, "-") && len(a) > 1 && a[1] != '-' {
			for _, c := range a[1:] {
				switch c {
				case 'l':
					long = true
				case 'a':
					all = true
				case 'A':
					almostAll = true
				case 'h':
					humanReadable = true
				case 'r':
					reverse = true
				case 't':
					sortByTime = true
				case 'S':
					sortBySize = true
				case 'U':
					noSort = true
				case 'R':
					recursive = true
				case 'i':
					inode = true
				case 'F':
					classify = true
				case '1':
					onePerLine = true
				case 'C':
					// columns - default
				case 'b':
					// escape - default
				}
			}
		} else if a == "--color" || a == "--color=always" || a == "--color=auto" {
			colorize = true
		} else if a == "--all" {
			all = true
		} else if a == "--almost-all" {
			almostAll = true
		} else if a == "--human-readable" {
			humanReadable = true
		} else if a == "--recursive" {
			recursive = true
		} else if a == "--inode" {
			inode = true
		} else if a == "--classify" {
			classify = true
		} else if a == "--reverse" {
			reverse = true
		} else if !strings.HasPrefix(a, "-") {
			dirs = append(dirs, a)
		}
	}

	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	// ANSI color codes
	colors := map[string]string{
		"dir":   "\033[01;34m",
		"link":  "\033[01;36m",
		"exec":  "\033[01;32m",
		"reset": "\033[0m",
	}
	if !colorize {
		for k := range colors {
			colors[k] = ""
		}
	}

	lookupName := func(uid uint32) string {
		u, err := user.LookupId(strconv.Itoa(int(uid)))
		if err != nil {
			return strconv.Itoa(int(uid))
		}
		return u.Username
	}
	lookupGroup := func(gid uint32) string {
		g, err := user.LookupGroupId(strconv.Itoa(int(gid)))
		if err != nil {
			return strconv.Itoa(int(gid))
		}
		return g.Name
	}

	printEntry := func(path string, info fs.FileInfo) {
		name := info.Name()
		suffix := ""
		if classify {
			switch {
			case info.Mode()&os.ModeSymlink != 0:
				suffix = "@"
			case info.IsDir():
				suffix = "/"
			case info.Mode()&0111 != 0:
				suffix = "*"
			}
		}
		if long {
			mode := info.Mode().String()
			nlink := uint64(1)
			uid := uint32(0)
			gid := uint32(0)
			ino := uint64(0)
			if stat, ok := info.Sys().(*syscall.Stat_t); ok {
				nlink = uint64(stat.Nlink)
				uid = stat.Uid
				gid = stat.Gid
				ino = stat.Ino
			}
			size := info.Size()
			sizeStr := strconv.FormatInt(size, 10)
			if humanReadable {
				sizeStr = humanSize(size)
			}
			modtime := info.ModTime().Format(time.Stamp)
			uname := lookupName(uid)
			gname := lookupGroup(gid)
			if inode {
				fmt.Printf("%8d ", ino)
			}
			colorStart, colorEnd := "", ""
			if colorize {
				if info.IsDir() {
					colorStart, colorEnd = colors["dir"], colors["reset"]
				} else if info.Mode()&os.ModeSymlink != 0 {
					colorStart, colorEnd = colors["link"], colors["reset"]
				} else if info.Mode()&0111 != 0 {
					colorStart, colorEnd = colors["exec"], colors["reset"]
				}
			}
			linkStr := ""
			if info.Mode()&os.ModeSymlink != 0 {
				target, _ := os.Readlink(path)
				linkStr = " -> " + target
			}
			fmt.Printf("%s %3d %-8s %-8s %8s %s %s%s%s%s%s\n",
				mode, nlink, uname, gname, sizeStr, modtime,
				colorStart, name+suffix, colorEnd, linkStr, "")
		} else {
			colorStart, colorEnd := "", ""
			if colorize {
				if info.IsDir() {
					colorStart, colorEnd = colors["dir"], colors["reset"]
				} else if info.Mode()&os.ModeSymlink != 0 {
					colorStart, colorEnd = colors["link"], colors["reset"]
				} else if info.Mode()&0111 != 0 {
					colorStart, colorEnd = colors["exec"], colors["reset"]
				}
			}
			if onePerLine || long {
				if inode {
					ino := uint64(0)
					if stat, ok := info.Sys().(*syscall.Stat_t); ok {
						ino = stat.Ino
					}
					fmt.Printf("%8d ", ino)
				}
				fmt.Printf("%s%s%s%s\n", colorStart, name, suffix, colorEnd)
			} else {
				fmt.Printf("%s%s%s%s  ", colorStart, name, suffix, colorEnd)
			}
		}
	}

	lsDir := func(dir string, header bool) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: %s: %v\n", dir, err)
			return
		}
		if header {
			fmt.Printf("%s:\n", dir)
		}

		type entry struct {
			name string
			info fs.FileInfo
			path string
		}
		var items []entry
		for _, e := range entries {
			name := e.Name()
			if !all && !almostAll && strings.HasPrefix(name, ".") {
				continue
			}
			if almostAll && (name == "." || name == "..") {
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			items = append(items, entry{name, info, filepath.Join(dir, name)})
		}
		if all {
			// Add . and ..
			di, _ := os.Lstat(dir)
			pi, _ := os.Lstat(filepath.Dir(dir))
			if di != nil {
				items = append([]entry{{".", di, dir}}, items...)
			}
			if pi != nil {
				items = append([]entry{{".", di, dir}, {"..", pi, filepath.Dir(dir)}}, items[1:]...)
			}
		}

		if !noSort {
			sort.SliceStable(items, func(i, j int) bool {
				a, b := items[i], items[j]
				var less bool
				switch {
				case sortByTime:
					less = a.info.ModTime().After(b.info.ModTime())
				case sortBySize:
					less = a.info.Size() > b.info.Size()
				default:
					less = strings.ToLower(a.name) < strings.ToLower(b.name)
				}
				if reverse {
					return !less
				}
				return less
			})
		}

		for _, item := range items {
			printEntry(item.path, item.info)
		}
		if !long && !onePerLine {
			fmt.Println()
		}

		if recursive {
			for _, item := range items {
				if item.info.IsDir() && item.name != "." && item.name != ".." {
					fmt.Println()
					lsDir(item.path, true)
				}
			}
		}
	}

	for i, d := range dirs {
		info, err := os.Lstat(d)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: %s: %v\n", d, err)
			continue
		}
		if info.IsDir() {
			if i > 0 {
				fmt.Println()
			}
			lsDir(d, len(dirs) > 1)
		} else {
			printEntry(d, info)
			if !long && !onePerLine {
				fmt.Println()
			}
		}
	}
}
