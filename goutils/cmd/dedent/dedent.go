// dedent - remove common leading whitespace from all lines (inverse of indent)
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	var lines []string
	for sc.Scan() { lines = append(lines, sc.Text()) }

	// Find minimum indent among non-empty lines
	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" { continue }
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent < 0 || indent < minIndent { minIndent = indent }
	}
	if minIndent < 0 { minIndent = 0 }

	for _, line := range lines {
		if len(line) >= minIndent { fmt.Println(line[minIndent:]) } else { fmt.Println(line) }
	}
}
