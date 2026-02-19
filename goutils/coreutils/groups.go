package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

func init() { register("groups", runGroups) }

func runGroups() {
	args := os.Args[1:]
	targets := []string{}
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			targets = append(targets, a)
		}
	}
	if len(targets) == 0 {
		targets = []string{""}
	}
	for _, name := range targets {
		var u *user.User
		var err error
		if name == "" {
			u, err = user.Current()
		} else {
			u, err = user.Lookup(name)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "groups: %v\n", err)
			continue
		}
		gids, err := u.GroupIds()
		if err != nil {
			fmt.Fprintf(os.Stderr, "groups: %v\n", err)
			continue
		}
		var gnames []string
		for _, gid := range gids {
			g, err := user.LookupGroupId(gid)
			if err != nil {
				gnames = append(gnames, strconv.Itoa(gid))
			} else {
				gnames = append(gnames, g.Name)
			}
		}
		if name != "" {
			fmt.Printf("%s : %s\n", name, strings.Join(gnames, " "))
		} else {
			fmt.Println(strings.Join(gnames, " "))
		}
	}
}
