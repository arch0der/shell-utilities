package main

import (
	"fmt"
	"os"
)

func init() { register("dircolors", runDircolors) }

const defaultLSColors = "rs=0:di=01;34:ln=01;36:mh=00:pi=40;33:so=01;35:do=01;35:bd=40;33;01:cd=40;33;01:or=40;31;01:mi=00:su=37;41:sg=30;43:ca=30;41:tw=30;42:ow=34;42:st=37;44:ex=01;32:*.tar=01;31:*.tgz=01;31:*.gz=01;31:*.zip=01;31"

const defaultDircolorsDB = `# Configuration file for dircolors
RESET 0
DIR 01;34
LINK 01;36
EXEC 01;32
`

func runDircolors() {
	args := os.Args[1:]
	shell := "sh"
	for _, a := range args {
		switch a {
		case "-b", "--sh", "--bourne-shell":
			shell = "sh"
		case "-c", "--csh", "--c-shell":
			shell = "csh"
		case "-p", "--print-database":
			fmt.Print(defaultDircolorsDB)
			return
		}
	}
	if shell == "csh" {
		fmt.Printf("setenv LS_COLORS '%s'\n", defaultLSColors)
	} else {
		fmt.Printf("LS_COLORS='%s';\nexport LS_COLORS\n", defaultLSColors)
	}
}
