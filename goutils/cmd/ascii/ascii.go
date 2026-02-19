// ascii - display ASCII table or look up characters/codepoints
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var ctrl = map[int]string{
	0:"NUL",1:"SOH",2:"STX",3:"ETX",4:"EOT",5:"ENQ",6:"ACK",
	7:"BEL",8:"BS",9:"TAB",10:"LF",11:"VT",12:"FF",13:"CR",
	14:"SO",15:"SI",16:"DLE",17:"DC1",18:"DC2",19:"DC3",20:"DC4",
	21:"NAK",22:"SYN",23:"ETB",24:"CAN",25:"EM",26:"SUB",27:"ESC",
	28:"FS",29:"GS",30:"RS",31:"US",127:"DEL",
}

func charLabel(i int) string {
	if n, ok := ctrl[i]; ok { return n }
	return string(rune(i))
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("%-5s %-6s %-6s  %s\n", "Dec", "Hex", "Oct", "Char")
		fmt.Println(strings.Repeat("-", 28))
		for i := 0; i < 128; i++ {
			fmt.Printf("%-5d %-6s %-6s  %s\n", i,
				fmt.Sprintf("0x%02X", i), fmt.Sprintf("0%03o", i), charLabel(i))
		}
		return
	}
	for _, arg := range os.Args[1:] {
		if n, err := strconv.ParseInt(arg, 0, 32); err == nil {
			r := rune(n)
			lbl := charLabel(int(n))
			if !unicode.IsPrint(r) { lbl = "<" + lbl + ">" }
			fmt.Printf("Dec:%-6d  Hex:0x%04X  Oct:0%04o  Char:%s\n", n, n, n, lbl)
		} else {
			for _, r := range arg {
				fmt.Printf("Dec:%-6d  Hex:0x%04X  Oct:0%04o  Char:%s\n", r, r, r, string(r))
			}
		}
	}
}
