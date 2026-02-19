// flip - flip text upside-down or mirror it horizontally
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// Unicode upside-down map
var flipMap = map[rune]rune{
	'a':'ɐ','b':'q','c':'ɔ','d':'p','e':'ǝ','f':'ɟ','g':'ƃ','h':'ɥ',
	'i':'ᴉ','j':'ɾ','k':'ʞ','l':'l','m':'ɯ','n':'u','o':'o','p':'d',
	'q':'b','r':'ɹ','s':'s','t':'ʇ','u':'n','v':'ʌ','w':'ʍ','x':'x',
	'y':'ʎ','z':'z',
	'A':'∀','B':'q','C':'Ɔ','D':'p','E':'Ǝ','F':'Ⅎ','G':'פ','H':'H',
	'I':'I','J':'ɾ','K':'ʞ','L':'˥','M':'W','N':'N','O':'O','P':'Ԁ',
	'Q':'Q','R':'ɹ','S':'S','T':'┴','U':'∩','V':'Λ','W':'M','X':'X',
	'Y':'⅄','Z':'Z',
	'0':'0','1':'Ɩ','2':'ᄅ','3':'Ɛ','4':'ㄣ','5':'ϛ','6':'9','7':'ㄥ',
	'8':'8','9':'6',
	'.':'˙',',':'\'','\'':',','!':'¡','?':'¿','(':')',')':'(',
	'[':']',']':'[','{':'}','}':'{','<':'>','>':'<','_':'‾','&':'⅋',
}

func flipUD(s string) string {
	// reverse the string and apply flip map
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 { runes[i], runes[j] = runes[j], runes[i] }
	for i, r := range runes { if m, ok := flipMap[r]; ok { runes[i] = m } }
	return string(runes)
}

func flipLR(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 { runes[i], runes[j] = runes[j], runes[i] }
	return string(runes)
}

func main() {
	mode := "ud"
	args := os.Args[1:]
	if len(args) > 0 && (args[0] == "-lr" || args[0] == "--mirror") { mode = "lr"; args = args[1:] }
	if len(args) > 0 && (args[0] == "-ud" || args[0] == "--flip") { mode = "ud"; args = args[1:] }

	fn := flipUD
	if mode == "lr" { fn = flipLR }

	if len(args) > 0 {
		text := strings.Join(args, " ")
		_ = utf8.RuneCountInString(text)
		fmt.Println(fn(text))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	var lines []string
	for sc.Scan() { lines = append(lines, fn(sc.Text())) }
	if mode == "ud" {
		for i := len(lines)-1; i >= 0; i-- { fmt.Println(lines[i]) }
	} else {
		for _, l := range lines { fmt.Println(l) }
	}
}
