// braille - convert text to Unicode Braille patterns
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Grade 1 Braille (alphabetic + digits)
var letterMap = map[rune]rune{
	'a':0x2801,'b':0x2803,'c':0x2809,'d':0x2819,'e':0x2811,
	'f':0x280B,'g':0x281B,'h':0x2813,'i':0x280A,'j':0x281A,
	'k':0x2805,'l':0x2807,'m':0x280D,'n':0x281D,'o':0x2815,
	'p':0x280F,'q':0x281F,'r':0x2817,'s':0x280E,'t':0x281E,
	'u':0x2825,'v':0x2827,'w':0x283A,'x':0x282D,'y':0x283D,
	'z':0x2835,
	'1':0x2801,'2':0x2803,'3':0x2809,'4':0x2819,'5':0x2811,
	'6':0x280B,'7':0x281B,'8':0x2813,'9':0x280A,'0':0x281A,
	' ':0x2800,',':0x2802,'.':0x2832,';':0x2806,':':0x2812,
	'?':0x2826,'!':0x2816,'-':0x2824,"'"[0]:0x2804,
}

func toBraille(s string) string {
	s = strings.ToLower(s)
	var sb strings.Builder
	for _, r := range s {
		if b, ok := letterMap[r]; ok {
			sb.WriteRune(b)
		} else if unicode.IsSpace(r) {
			sb.WriteRune(0x2800)
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) > 1 {
		fmt.Println(toBraille(strings.Join(os.Args[1:], " ")))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { fmt.Println(toBraille(sc.Text())) }
}
