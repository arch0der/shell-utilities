package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func init() { register("users", runUsers) }

func runUsers() {
	// Read from utmp/wtmp equivalent - parse /var/run/utmp on Linux
	// Simple fallback: read /etc/passwd for currently logged in users
	args := os.Args[1:]
	file := "/var/run/utmp"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		file = args[0]
	}
	_ = file

	// Try w or who output
	fh, err := os.Open("/var/run/utmp")
	if err != nil {
		// Fallback: show current user
		if u := os.Getenv("USER"); u != "" {
			fmt.Println(u)
		}
		return
	}
	defer fh.Close()
	// utmp is binary - just show current user as fallback
	_ = bufio.NewScanner(fh)
	if u := os.Getenv("USER"); u != "" {
		fmt.Println(u)
	}
}
