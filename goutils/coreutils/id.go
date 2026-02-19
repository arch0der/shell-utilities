package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

func init() { register("id", runId) }

func runId() {
	args := os.Args[1:]
	printUser := false
	printGroup := false
	printGroups := false
	nameOnly := false
	realID := false
	zero := false
	targets := []string{}

	for _, a := range args {
		switch a {
		case "-u", "--user":
			printUser = true
		case "-g", "--group":
			printGroup = true
		case "-G", "--groups":
			printGroups = true
		case "-n", "--name":
			nameOnly = true
		case "-r", "--real":
			realID = true
		case "-z", "--zero":
			zero = true
		default:
			if !strings.HasPrefix(a, "-") {
				targets = append(targets, a)
			}
		}
	}
	_ = realID

	var u *user.User
	var err error
	if len(targets) > 0 {
		u, err = user.Lookup(targets[0])
	} else {
		u, err = user.Current()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "id: %v\n", err)
		os.Exit(1)
	}

	sep := "\n"
	if zero {
		sep = "\x00"
	}

	if printUser {
		if nameOnly {
			fmt.Print(u.Username + sep)
		} else {
			fmt.Print(u.Uid + sep)
		}
		return
	}
	if printGroup {
		g, _ := user.LookupGroupId(u.Gid)
		if nameOnly && g != nil {
			fmt.Print(g.Name + sep)
		} else {
			fmt.Print(u.Gid + sep)
		}
		return
	}
	if printGroups {
		gids, _ := u.GroupIds()
		var parts []string
		for _, gid := range gids {
			if nameOnly {
				g, err := user.LookupGroupId(gid)
				if err == nil {
					parts = append(parts, g.Name)
					continue
				}
			}
			parts = append(parts, strconv.Itoa(gid))
		}
		if zero {
			fmt.Print(strings.Join(parts, sep) + sep)
		} else {
			fmt.Println(strings.Join(parts, " "))
		}
		return
	}

	// Full id output
	g, _ := user.LookupGroupId(u.Gid)
	gname := u.Gid
	if g != nil {
		gname = g.Name
	}
	fmt.Printf("uid=%s(%s) gid=%s(%s)", u.Uid, u.Username, u.Gid, gname)
	gids, _ := u.GroupIds()
	if len(gids) > 0 {
		fmt.Print(" groups=")
		var groupParts []string
		for _, gid := range gids {
			grp, err := user.LookupGroupId(gid)
			if err != nil {
				groupParts = append(groupParts, gid+"(unknown)")
			} else {
				groupParts = append(groupParts, gid+"("+grp.Name+")")
			}
		}
		fmt.Print(strings.Join(groupParts, ","))
	}
	fmt.Println()
}
