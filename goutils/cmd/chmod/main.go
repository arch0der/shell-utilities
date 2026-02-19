// chmod - Change file permissions
// Usage: chmod [-R] mode file...
// Mode: octal (755) or symbolic (u+x, go-w, a=r, etc.)
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var recursive = flag.Bool("R", false, "Recursively change permissions")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: chmod [-R] mode file...")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	modeStr := flag.Arg(0)
	exitCode := 0

	for _, path := range flag.Args()[1:] {
		if err := chmodPath(path, modeStr); err != nil {
			fmt.Fprintln(os.Stderr, "chmod:", err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func chmodPath(path, modeStr string) error {
	if *recursive {
		return filepath.Walk(path, func(p string, _ os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			return applyMode(p, modeStr)
		})
	}
	return applyMode(path, modeStr)
}

func applyMode(path, modeStr string) error {
	// Octal mode
	if n, err := strconv.ParseUint(modeStr, 8, 32); err == nil {
		return os.Chmod(path, os.FileMode(n))
	}

	// Symbolic mode: [ugoa][+-=][rwxst...]
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	mode := info.Mode().Perm()

	for _, clause := range strings.Split(modeStr, ",") {
		clause = strings.TrimSpace(clause)
		who, op, perms := parseClause(clause)
		var bits os.FileMode
		for _, p := range perms {
			switch p {
			case 'r':
				bits |= 0444
			case 'w':
				bits |= 0222
			case 'x':
				bits |= 0111
			}
		}
		// Mask by who
		var mask os.FileMode
		for _, w := range who {
			switch w {
			case 'u':
				mask |= 0700
			case 'g':
				mask |= 0070
			case 'o':
				mask |= 0007
			case 'a':
				mask = 0777
			}
		}
		if mask == 0 {
			mask = 0777
		}
		switch op {
		case '+':
			mode |= bits & mask
		case '-':
			mode &^= bits & mask
		case '=':
			mode = (mode &^ mask) | (bits & mask)
		}
	}
	return os.Chmod(path, mode)
}

func parseClause(s string) (who string, op byte, perms string) {
	i := 0
	for i < len(s) && (s[i] == 'u' || s[i] == 'g' || s[i] == 'o' || s[i] == 'a') {
		who += string(s[i])
		i++
	}
	if i < len(s) && (s[i] == '+' || s[i] == '-' || s[i] == '=') {
		op = s[i]
		i++
	}
	perms = s[i:]
	return
}
