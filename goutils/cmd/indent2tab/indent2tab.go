// indent2tab - convert space-indented code to tab-indented (or vice versa)
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	toSpaces := false
	spaceWidth := 4
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-s", "--spaces": toSpaces = true
		case "-w": i++; spaceWidth, _ = strconv.Atoi(args[i])
		}
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		if toSpaces {
			// tabs → spaces
			var b strings.Builder
			for _, ch := range line {
				if ch == '\t' { b.WriteString(strings.Repeat(" ", spaceWidth)) } else { b.WriteRune(ch); break }
			}
			// find where indent ends
			rest := strings.TrimLeft(line, "\t")
			tabs := len(line) - len(rest)
			fmt.Println(strings.Repeat(" ", tabs*spaceWidth) + rest)
		} else {
			// spaces → tabs
			rest := strings.TrimLeft(line, " ")
			spaces := len(line) - len(rest)
			tabs := spaces / spaceWidth
			rem := spaces % spaceWidth
			fmt.Print(strings.Repeat("\t", tabs) + strings.Repeat(" ", rem) + rest + "\n")
		}
	}
}
