package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

func init() { register("sync", runSync) }

func runSync() {
	args := os.Args[1:]
	data := false
	files := []string{}
	for _, a := range args {
		if a == "--data" {
			data = true
		} else if !strings.HasPrefix(a, "-") {
			files = append(files, a)
		}
	}
	if len(files) > 0 {
		for _, f := range files {
			fh, err := os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "sync: %s: %v\n", f, err)
				continue
			}
			if data {
				fh.Sync()
			} else {
				fh.Sync()
			}
			fh.Close()
		}
		return
	}
	syscall.Sync()
}
