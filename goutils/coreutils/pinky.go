package main

import (
	"fmt"
	"os"
	"os/user"
)

func init() { register("pinky", runPinky) }

func runPinky() {
	args := os.Args[1:]
	short := false
	noHeader := false
	targets := []string{}

	for _, a := range args {
		switch a {
		case "-s":
			short = true
		case "-f":
			noHeader = true
		default:
			if len(a) > 0 && a[0] != '-' {
				targets = append(targets, a)
			}
		}
	}

	if !noHeader && !short {
		fmt.Printf("%-10s %-20s %-15s %s\n", "Login", "Name", "TTY", "Idle")
	}

	printUser := func(u *user.User) {
		name := u.Name
		if name == "" {
			name = u.Username
		}
		if short {
			fmt.Printf("%-10s %-20s\n", u.Username, name)
		} else {
			fmt.Printf("%-10s %-20s %-15s %s\n", u.Username, name, "?", "?")
		}
	}

	if len(targets) == 0 {
		u, err := user.Current()
		if err == nil {
			printUser(u)
		}
		return
	}
	for _, t := range targets {
		u, err := user.Lookup(t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pinky: %s: no such user\n", t)
			continue
		}
		printUser(u)
	}
}
