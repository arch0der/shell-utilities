// cp - Copy files or directories
// Usage: cp [-r] <src> <dst>
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var recursive = flag.Bool("r", false, "Copy directories recursively")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cp [-r] <src> <dst>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}
	src, dst := flag.Arg(0), flag.Arg(1)
	if err := copyPath(src, dst); err != nil {
		fmt.Fprintln(os.Stderr, "cp:", err)
		os.Exit(1)
	}
}

func copyPath(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		if !*recursive {
			return fmt.Errorf("%s is a directory (use -r)", src)
		}
		return copyDir(src, dst, info)
	}
	// If dst is an existing directory, copy file inside it
	if di, err2 := os.Stat(dst); err2 == nil && di.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}
	return copyFile(src, dst, info.Mode())
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func copyDir(src, dst string, info os.FileInfo) error {
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if err := copyPath(s, d); err != nil {
			return err
		}
	}
	return nil
}
