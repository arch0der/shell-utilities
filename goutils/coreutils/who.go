package main

import (
	"fmt"
	"os"
	"os/user"
	"time"
)

func init() {
	register("who", runWho)
	register("whoami", runWhoami)
	register("logname", runLogname)
}

func runWho() {
	args := os.Args[1:]
	amI := false
	for _, a := range args {
		if a == "-m" || a == "am" || a == "i" {
			amI = true
		}
	}
	u, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "who:", err)
		os.Exit(1)
	}
	tty := os.Getenv("TTY")
	if tty == "" {
		tty = "tty1"
	}
	now := time.Now().Format("2006-01-02 15:04")
	if amI {
		h, _ := os.Hostname()
		fmt.Printf("%s\t%s\t%s (%s)\n", u.Username, tty, now, h)
	} else {
		fmt.Printf("%-12s %-8s %s\n", u.Username, tty, now)
	}
}

func runWhoami() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "whoami:", err)
		os.Exit(1)
	}
	fmt.Println(u.Username)
}

func runLogname() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "logname:", err)
		os.Exit(1)
	}
	fmt.Println(u.Username)
}
