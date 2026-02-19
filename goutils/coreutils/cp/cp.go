package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	recursive := false
	force := false
	preserve := false
	interactive := false
	verbose := false
	files := []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-r" || a == "-R" || a == "--recursive":
			recursive = true
		case a == "-f" || a == "--force":
			force = true
		case a == "-p" || a == "--preserve":
			preserve = true
		case a == "-i" || a == "--interactive":
			interactive = true
		case a == "-v" || a == "--verbose":
			verbose = true
		case a == "-a" || a == "--archive":
			recursive, preserve = true, true
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}
	_ = force
	_ = interactive
	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "cp: missing destination")
		os.Exit(1)
	}
	dest := files[len(files)-1]
	srcs := files[:len(files)-1]

	copyFile := func(src, dst string) error {
		srcInfo, err := os.Stat(src)
		if err != nil {
			return err
		}
		if srcInfo.IsDir() {
			if !recursive {
				return fmt.Errorf("omitting directory '%s'", src)
			}
			os.MkdirAll(dst, srcInfo.Mode())
			entries, _ := os.ReadDir(src)
			for _, e := range entries {
				copyFile(filepath.Join(src, e.Name()), filepath.Join(dst, e.Name()))
			}
			return nil
		}
		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()
		mode := srcInfo.Mode()
		out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
		if err != nil {
			return err
		}
		defer out.Close()
		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		if preserve {
			os.Chtimes(dst, srcInfo.ModTime(), srcInfo.ModTime())
		}
		return nil
	}

	destInfo, destErr := os.Stat(dest)
	exitCode := 0
	for _, src := range srcs {
		dst := dest
		if destErr == nil && destInfo.IsDir() {
			dst = filepath.Join(dest, filepath.Base(src))
		}
		if verbose {
			fmt.Printf("'%s' -> '%s'\n", src, dst)
		}
		if err := copyFile(src, dst); err != nil {
			fmt.Fprintf(os.Stderr, "cp: %v\n", err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
