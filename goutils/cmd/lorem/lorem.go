// lorem - generate Lorem Ipsum placeholder text
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var words = strings.Fields(`lorem ipsum dolor sit amet consectetur adipiscing elit
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua ut enim ad minim
veniam quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
consequat duis aute irure dolor in reprehenderit in voluptate velit esse cillum
dolore eu fugiat nulla pariatur excepteur sint occaecat cupidatat non proident
sunt in culpa qui officia deserunt mollit anim id est laborum curabitur pretium
tincidunt lacus nulla aliquet enim cras eget est lorem ipsum dolor sit amet
consectetur adipiscing elit proin blandit erat consectetur lorem pretium dignissim`)

func sentence(rng *rand.Rand, minW, maxW int) string {
	n := minW + rng.Intn(maxW-minW+1)
	ws := make([]string, n)
	for i := range ws { ws[i] = words[rng.Intn(len(words))] }
	s := strings.Join(ws, " ")
	if len(s) > 0 { s = strings.ToUpper(s[:1]) + s[1:] }
	return s + "."
}

func paragraph(rng *rand.Rand, sentences int) string {
	ss := make([]string, sentences)
	for i := range ss { ss[i] = sentence(rng, 8, 18) }
	return strings.Join(ss, " ")
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	mode := "paragraphs"
	n := 3
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-w", "--words": mode = "words"; if i+1 < len(args) { n, _ = strconv.Atoi(args[i+1]); i++ }
		case "-s", "--sentences": mode = "sentences"; if i+1 < len(args) { n, _ = strconv.Atoi(args[i+1]); i++ }
		case "-p", "--paragraphs": mode = "paragraphs"; if i+1 < len(args) { n, _ = strconv.Atoi(args[i+1]); i++ }
		default:
			if v, err := strconv.Atoi(args[i]); err == nil { n = v }
		}
	}

	switch mode {
	case "words":
		ws := make([]string, n)
		for i := range ws { ws[i] = words[rng.Intn(len(words))] }
		fmt.Println(strings.Join(ws, " "))
	case "sentences":
		ss := make([]string, n)
		for i := range ss { ss[i] = sentence(rng, 8, 18) }
		fmt.Println(strings.Join(ss, " "))
	case "paragraphs":
		for i := 0; i < n; i++ {
			fmt.Println(paragraph(rng, 4+rng.Intn(4)))
			if i < n-1 { fmt.Println() }
		}
	}
}
