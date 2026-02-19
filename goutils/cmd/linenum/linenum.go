// linenum - add or remove line numbers from text
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	remove := false
	start := 1
	step := 1
	format := "%d\t"
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-r", "--remove": remove = true
		case "-s": i++; start, _ = strconv.Atoi(args[i])
		case "--step": i++; step, _ = strconv.Atoi(args[i])
		case "-f": i++; format = args[i]
		}
	}

	sc := bufio.NewScanner(os.Stdin)
	lineN := start
	if remove {
		re := regexp.MustCompile(`^\s*\d+\s*[\t:]?\s?`)
		for sc.Scan() { fmt.Println(re.ReplaceAllString(sc.Text(), "")) }
		return
	}
	for sc.Scan() {
		line := sc.Text()
		num := strings.ReplaceAll(format, "%d", strconv.Itoa(lineN))
		fmt.Print(num + line + "\n")
		lineN += step
	}
}
