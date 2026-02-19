// diceware - generate diceware passphrases using cryptographic randomness
package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// 100-word subset of EFF large wordlist (real diceware uses 7776 words)
var words = []string{
	"abacus","abbey","abbot","abduct","abide","abiotic","ablaze","able","abode","abort",
	"about","above","abrupt","absent","absorb","abstract","absurd","accent","access","account",
	"acidify","acorn","acre","across","act","action","active","actor","acute","adapt",
	"address","adept","adjust","adopt","advance","advice","afford","afraid","after","agenda",
	"agree","alarm","album","alert","algae","align","alley","allow","ally","almond",
	"alone","aloof","alpine","alter","altimeter","alumni","amber","amble","amend","ample",
	"analog","ancient","angle","annex","answer","antler","anvil","apart","apex","apple",
	"apricot","arch","arctic","area","argue","arid","armband","armor","arrow","asking",
	"aspect","asset","atlas","atom","attic","autumn","avid","awake","award","awful",
	"axle","azure","bacon","badge","bagel","bamboo","banana","barley","barn","basket",
}

func rollDice() int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
	return int(n.Int64())
}

func main() {
	wordCount := 6
	count := 1
	sep := "-"
	entropy := false
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-n": i++; count, _ = strconv.Atoi(args[i])
		case "-s": i++; sep = args[i]
		case "-e": entropy = true
		default:
			if n, err := strconv.Atoi(args[i]); err == nil { wordCount = n }
		}
	}

	for i := 0; i < count; i++ {
		ws := make([]string, wordCount)
		for j := range ws { ws[j] = words[rollDice()] }
		fmt.Println(strings.Join(ws, sep))
	}

	if entropy {
		bits := float64(wordCount) * 6.89 // log2(100)
		fmt.Fprintf(os.Stderr, "\nEntropy: ~%.0f bits (%d words from %d-word list)\n", bits, wordCount, len(words))
		fmt.Fprintln(os.Stderr, "Note: Use a full 7776-word EFF list for production use.")
	}
}
