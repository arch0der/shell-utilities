// fingerprint - generate visual fingerprints (identicons) for hashes/strings
package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
)

// SSH randomart-style visual fingerprint
func randomart(data []byte, label string) string {
	const (
		w, h = 17, 9
		chars = " .o+=*BOX@%&#/^SE"
	)
	field := make([]int, w*h)
	x, y := w/2, h/2
	for _, b := range data {
		for bit := 0; bit < 4; bit++ {
			dx := (int(b>>(bit*2+1))&1)*2 - 1
			dy := (int(b>>(bit*2))&1)*2 - 1
			x += dx; y += dy
			if x < 0 { x = 0 }; if x >= w { x = w-1 }
			if y < 0 { y = 0 }; if y >= h { y = h-1 }
			if field[y*w+x] < len(chars)-3 { field[y*w+x]++ }
		}
	}
	field[h/2*w+w/2] = len(chars) - 1 // S = start
	field[(h-1)*w+(w-1)] = len(chars) - 2 // E = end ... simplified

	var sb strings.Builder
	topLabel := fmt.Sprintf("[%s]", label)
	pad := (w+2-len(topLabel))/2
	sb.WriteString("+" + strings.Repeat("-", pad) + topLabel + strings.Repeat("-", w+2-pad-len(topLabel)) + "+\n")
	for row := 0; row < h; row++ {
		sb.WriteByte('|')
		for col := 0; col < w; col++ { sb.WriteByte(chars[field[row*w+col]]) }
		sb.WriteString("|\n")
	}
	sb.WriteString("+" + strings.Repeat("-", w) + "+")
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: fingerprint <string_or_file>")
		os.Exit(1)
	}
	input := strings.Join(os.Args[1:], " ")
	label := input
	if len(label) > 12 { label = label[:12] }

	// Try as file
	if data, err := os.ReadFile(input); err == nil {
		h := sha256.Sum256(data)
		fmt.Printf("SHA256: %x\n", h)
		fmt.Println(randomart(h[:], "SHA256"))
		return
	}

	h := sha256.Sum256([]byte(input))
	fmt.Printf("Input : %s\n", input)
	fmt.Printf("SHA256: %x\n", h)
	fmt.Println(randomart(h[:], label))
}
