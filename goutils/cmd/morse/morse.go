// morse - encode/decode Morse code
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var toMorse = map[rune]string{
	'A':".-",'B':"-...",'C':"-.-.",'D':"-..",'E':".",'F':"..-.",'G':"--.",
	'H':"....",'I':"..",'J':".---",'K':"-.-",'L':".-..",'M':"--",'N':"-.",
	'O':"---",'P':".--.",'Q':"--.-",'R':".-.",'S':"...",'T':"-",'U':"..-",
	'V':"...-",'W':".--",'X':"-..-",'Y':"-.--",'Z':"--..",
	'0':"-----",'1':".----",'2':"..---",'3':"...--",'4':"....-",'5':".....",
	'6':"-....","7":"--...", // fix: use rune keys
	'7':"--...",'8':"---..",'9':"----.",
	'.':".-.-.-",',':"--..--",'?':"..--..",'!':"-.-.--",'/':"-..-.","=":"-...-",
	'+':".-.-.",'-':"-....-",'_':"..--.-",'"':".-..-.",'$':"...-..-",'&':".-...",
	' ':"/",
}

var fromMorse map[string]rune

func init() {
	fromMorse = map[string]rune{}
	for k, v := range toMorse { if k != ' ' { fromMorse[v] = k } }
	fromMorse["/"] = ' '
}

func encode(text string) string {
	text = strings.ToUpper(text)
	var parts []string
	for _, r := range text {
		if m, ok := toMorse[r]; ok { parts = append(parts, m) }
	}
	return strings.Join(parts, " ")
}

func decode(text string) string {
	var sb strings.Builder
	words := strings.Split(text, " / ")
	for wi, word := range words {
		letters := strings.Fields(word)
		for _, l := range letters {
			if r, ok := fromMorse[l]; ok { sb.WriteRune(r) } else { sb.WriteString("?") }
		}
		if wi < len(words)-1 { sb.WriteRune(' ') }
	}
	return sb.String()
}

func main() {
	mode := "encode"
	if len(os.Args) > 1 && (os.Args[1] == "-d" || os.Args[1] == "decode") { mode = "decode" }

	process := func(line string) string {
		if mode == "encode" { return encode(line) }
		return decode(line)
	}

	args := os.Args[1:]
	if len(args) > 0 && (args[0] == "-d" || args[0] == "decode" || args[0] == "encode") { args = args[1:] }
	if len(args) > 0 {
		fmt.Println(process(strings.Join(args, " ")))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { fmt.Println(process(sc.Text())) }
}
