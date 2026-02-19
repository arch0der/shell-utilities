// phonetic - convert letters to NATO phonetic alphabet
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var nato = map[rune]string{
	'A':"Alpha",'B':"Bravo",'C':"Charlie",'D':"Delta",'E':"Echo",
	'F':"Foxtrot",'G':"Golf",'H':"Hotel",'I':"India",'J':"Juliet",
	'K':"Kilo",'L':"Lima",'M':"Mike",'N':"November",'O':"Oscar",
	'P':"Papa",'Q':"Quebec",'R':"Romeo",'S':"Sierra",'T':"Tango",
	'U':"Uniform",'V':"Victor",'W':"Whiskey",'X':"X-ray",'Y':"Yankee",
	'Z':"Zulu",
	'0':"Zero",'1':"One",'2':"Two",'3':"Three",'4':"Four",
	'5':"Five",'6':"Six",'7':"Seven",'8':"Eight",'9':"Niner",
	'.':"Period",',':"Comma",'-':"Dash",'_':"Underscore",
	'@':"At",'#':"Hash",'/':"Slash",'!':"Exclamation",'?':"Question",
}

func convert(s string) string {
	var parts []string
	for _, r := range strings.ToUpper(s) {
		if r == ' ' { parts = append(parts, " / "); continue }
		if !unicode.IsPrint(r) { continue }
		if word, ok := nato[r]; ok { parts = append(parts, word) } else { parts = append(parts, string(r)) }
	}
	return strings.Join(parts, " ")
}

func main() {
	if len(os.Args) > 1 {
		fmt.Println(convert(strings.Join(os.Args[1:], " ")))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { fmt.Println(convert(sc.Text())) }
}
