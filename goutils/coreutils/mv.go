package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func init() { register("mv", runMv) }

func runMv() {
	args := os.Args[1:]
	force := false
	interactive := false
	verbose := false
	noOverwrite := false
	backup := false
	files := []string{}

	for _, a := range args {
		switch a {
		case "-f", "--force":
			force = true
		case "-i", "--interactive":
			interactive = true
		case "-v", "--verbose":
			verbose = true
		case "-n", "--no-clobber":
			noOverwrite = true
		case "-b", "--backup":
			backup = true
		default:
			if !strings.HasPrefix(a, "-") {
				files = append(files, a)
			}
		}
	}
	_ = force
	_ = interactive

	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "mv: missing destination")
		os.Exit(1)
	}

	dest := files[len(files)-1]
	srcs := files[:len(files)-1]

	exitCode := 0
	for _, src := range srcs {
		dst := dest
		if info, err := os.Stat(dest); err == nil && info.IsDir() {
			dst = filepath.Join(dest, filepath.Base(src))
		}
		if noOverwrite {
			if _, err := os.Stat(dst); err == nil {
				continue
			}
		}
		if backup {
			if _, err := os.Stat(dst); err == nil {
				os.Rename(dst, dst+"~")
			}
		}
		err := os.Rename(src, dst)
		if err != nil {
			// Cross-device: copy then delete
			if copyErr := crossDeviceCopy(src, dst); copyErr != nil {
				fmt.Fprintf(os.Stderr, "mv: %v\n", copyErr)
				exitCode = 1
				continue
			}
			os.RemoveAll(src)
		}
		if verbose {
			fmt.Printf("'%s' -> '%s'\n", src, dst)
		}
	}
	os.Exit(exitCode)
}

func crossDeviceCopy(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if srcInfo.IsDir() {
		if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
			return err
		}
		entries, _ := os.ReadDir(src)
		for _, e := range entries {
			crossDeviceCopy(filepath.Join(src, e.Name()), filepath.Join(dst, e.Name()))
		}
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
