// install - Copy files with permissions and ownership
// Usage: install [-d] [-m mode] [-o owner] [-g group] source... dest
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

var (
	mode      = flag.String("m", "755", "Permission mode")
	owner     = flag.String("o", "", "Owner")
	group     = flag.String("g", "", "Group")
	mkdirFlag = flag.Bool("d", false, "Create directories")
	verbose   = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: install [-d] [-m mode] [-o owner] [-g group] source... dest")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	perm, err := strconv.ParseUint(*mode, 8, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, "install: invalid mode:", *mode)
		os.Exit(1)
	}
	fileMode := os.FileMode(perm)

	uid, gid := -1, -1
	if *owner != "" {
		u, err := user.Lookup(*owner)
		if err != nil {
			fmt.Fprintln(os.Stderr, "install: invalid owner:", *owner)
			os.Exit(1)
		}
		uid, _ = strconv.Atoi(u.Uid)
	}
	if *group != "" {
		g, err := user.LookupGroup(*group)
		if err != nil {
			fmt.Fprintln(os.Stderr, "install: invalid group:", *group)
			os.Exit(1)
		}
		gid, _ = strconv.Atoi(g.Gid)
	}

	if *mkdirFlag {
		// Create directories
		exitCode := 0
		for _, dir := range flag.Args() {
			if err := os.MkdirAll(dir, fileMode); err != nil {
				fmt.Fprintln(os.Stderr, "install:", err)
				exitCode = 1
				continue
			}
			os.Chmod(dir, fileMode)
			if uid >= 0 || gid >= 0 {
				syscall.Chown(dir, uid, gid)
			}
			if *verbose {
				fmt.Printf("install: creating directory '%s'\n", dir)
			}
		}
		os.Exit(exitCode)
	}

	// Copy files
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	dest := flag.Arg(flag.NArg() - 1)
	sources := flag.Args()[:flag.NArg()-1]

	destIsDir := false
	if info, err := os.Stat(dest); err == nil && info.IsDir() {
		destIsDir = true
	}

	exitCode := 0
	for _, src := range sources {
		var target string
		if destIsDir {
			target = dest + "/" + fileNameFromPath(src)
		} else {
			target = dest
		}

		if err := installFile(src, target, fileMode, uid, gid); err != nil {
			fmt.Fprintln(os.Stderr, "install:", err)
			exitCode = 1
			continue
		}
		if *verbose {
			fmt.Printf("'%s' -> '%s'\n", src, target)
		}
	}
	os.Exit(exitCode)
}

func installFile(src, dst string, mode os.FileMode, uid, gid int) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer df.Close()

	if _, err := io.Copy(df, sf); err != nil {
		return err
	}
	df.Close()

	if err := os.Chmod(dst, mode); err != nil {
		return err
	}
	if uid >= 0 || gid >= 0 {
		syscall.Chown(dst, uid, gid)
	}
	return nil
}

func fileNameFromPath(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
