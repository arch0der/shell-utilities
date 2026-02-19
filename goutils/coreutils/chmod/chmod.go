package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	recursive := false
	modeStr := ""
	files := []string{}
	for _, a := range args {
		if a == "-R" || a == "-r" || a == "--recursive" {
			recursive = true
		} else if !strings.HasPrefix(a, "-") {
			if modeStr == "" {
				modeStr = a
			} else {
				files = append(files, a)
			}
		}
	}
	if modeStr == "" || len(files) == 0 {
		fmt.Fprintln(os.Stderr, "chmod: missing operand")
		os.Exit(1)
	}
	parseMode := func(path string, modeStr string) (os.FileMode, error) {
		info, err := os.Lstat(path)
		if err != nil {
			return 0, err
		}
		current := info.Mode()
		// Numeric mode
		if n, err := strconv.ParseUint(modeStr, 8, 32); err == nil {
			return os.FileMode(n), nil
		}
		// Symbolic mode: [ugoa][+-=][rwxXst]
		mode := current
		for _, part := range strings.Split(modeStr, ",") {
			who := ""
			op := byte(0)
			perms := ""
			for i, c := range part {
				if c == '+' || c == '-' || c == '=' {
					who = part[:i]
					op = byte(c)
					perms = part[i+1:]
					break
				}
			}
			if op == 0 {
				continue
			}
			if who == "" {
				who = "a"
			}
			var bits os.FileMode
			for _, p := range perms {
				switch p {
				case 'r':
					bits |= 0444
				case 'w':
					bits |= 0222
				case 'x':
					bits |= 0111
				case 'X':
					if current.IsDir() || (current&0111 != 0) {
						bits |= 0111
					}
				case 's':
					bits |= os.ModeSetuid | os.ModeSetgid
				case 't':
					bits |= os.ModeSticky
				}
			}
			// filter by who
			var mask os.FileMode = 0777
			switch who {
			case "u":
				mask = 0700
				bits &= 0700
			case "g":
				mask = 0070
				bits &= 0070
			case "o":
				mask = 0007
				bits &= 0007
			}
			switch op {
			case '+':
				mode |= bits
			case '-':
				mode &^= bits
			case '=':
				mode = (mode &^ mask) | bits
			}
		}
		return mode, nil
	}
	var doChmod func(path string)
	doChmod = func(path string) {
		m, err := parseMode(path, modeStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "chmod: %s: %v\n", path, err)
			return
		}
		if err := os.Chmod(path, m); err != nil {
			fmt.Fprintf(os.Stderr, "chmod: %s: %v\n", path, err)
		}
		if recursive {
			info, err := os.Lstat(path)
			if err == nil && info.IsDir() {
				entries, _ := os.ReadDir(path)
				for _, e := range entries {
					doChmod(filepath.Join(path, e.Name()))
				}
			}
		}
	}
	for _, f := range files {
		doChmod(f)
	}
}
