package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func init() { register("touch", runTouch) }

func runTouch() {
	args := os.Args[1:]
	noCreate := false
	accessOnly := false
	modifyOnly := false
	timeStr := ""
	files := []string{}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-c" || a == "--no-create":
			noCreate = true
		case a == "-a":
			accessOnly = true
		case a == "-m":
			modifyOnly = true
		case a == "-t" && i+1 < len(args):
			i++
			timeStr = args[i]
		case a == "-d" && i+1 < len(args):
			i++
			timeStr = args[i]
		case !strings.HasPrefix(a, "-"):
			files = append(files, a)
		}
	}
	_ = accessOnly
	_ = modifyOnly

	var t time.Time
	if timeStr != "" {
		layouts := []string{"200601021504.05", "200601021504", "2006-01-02 15:04:05", "2006-01-02"}
		for _, l := range layouts {
			var err error
			t, err = time.Parse(l, timeStr)
			if err == nil {
				break
			}
		}
		if t.IsZero() {
			t = time.Now()
		}
	} else {
		t = time.Now()
	}

	exitCode := 0
	for _, f := range files {
		_, err := os.Stat(f)
		if os.IsNotExist(err) {
			if noCreate {
				continue
			}
			fh, err := os.Create(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "touch: %s: %v\n", f, err)
				exitCode = 1
				continue
			}
			fh.Close()
		}
		if err := os.Chtimes(f, t, t); err != nil {
			fmt.Fprintf(os.Stderr, "touch: %s: %v\n", f, err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
