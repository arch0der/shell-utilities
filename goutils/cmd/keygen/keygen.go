// keygen - generate cryptographic keys, tokens, and secrets
package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

const (
	alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	urlSafe      = alphanumeric + "-_"
	humanReadable = "abcdefghjkmnpqrstuvwxyz23456789" // no ambiguous chars
)

func randBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func randString(charset string, n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		sb.WriteByte(charset[idx.Int64()])
	}
	return sb.String()
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: keygen <type> [length]
  hex      [bytes=32]     random hex string
  base64   [bytes=32]     random base64 string
  token    [chars=32]     URL-safe alphanumeric token
  secret   [chars=64]     long URL-safe secret
  human    [chars=20]     human-readable (no ambiguous chars)
  uuid                    random UUIDv4
  pin      [digits=6]     numeric PIN
  passphrase [words=4]    word-list passphrase (EFF-style subset)`)
	os.Exit(1)
}

var wordList = []string{
	"apple","brave","crane","dance","eagle","flame","grape","horse","image","jolly",
	"knife","lemon","mango","noble","ocean","piano","queen","river","stone","tiger",
	"uncle","vivid","witch","xenon","yacht","zebra","amber","blaze","clown","drove",
	"ember","frost","glare","haste","ivory","joker","knack","lunar","maple","nerve",
	"onion","pixie","quirk","radar","snowy","thick","ultra","vapor","wider","xerox",
	"yearn","zonal","audio","block","crisp","delta","equip","fancy","gnome","happy",
	"index","jelly","karma","light","month","north","opera","plume","quote","rally",
	"sigma","tower","umbra","value","world","exact","young","zilch","arrow","bring",
	"cubic","depot","evoke","fancy","globe","honey","inner","jazzy","kudos","lived",
	"magic","night","orbit","proxy","quiet","rocky","solar","tidal","unfold","visor",
}

func main() {
	if len(os.Args) < 2 { usage() }
	typ := os.Args[1]
	n := 0
	if len(os.Args) > 2 { n, _ = strconv.Atoi(os.Args[2]) }

	switch typ {
	case "hex":
		if n == 0 { n = 32 }
		fmt.Println(hex.EncodeToString(randBytes(n)))
	case "base64":
		if n == 0 { n = 32 }
		fmt.Println(base64.URLEncoding.EncodeToString(randBytes(n)))
	case "token":
		if n == 0 { n = 32 }
		fmt.Println(randString(urlSafe, n))
	case "secret":
		if n == 0 { n = 64 }
		fmt.Println(randString(urlSafe, n))
	case "human":
		if n == 0 { n = 20 }
		fmt.Println(randString(humanReadable, n))
	case "uuid":
		b := randBytes(16)
		b[6] = (b[6] & 0x0f) | 0x40
		b[8] = (b[8] & 0x3f) | 0x80
		fmt.Printf("%x-%x-%x-%x-%x\n", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	case "pin":
		if n == 0 { n = 6 }
		fmt.Println(randString("0123456789", n))
	case "passphrase":
		if n == 0 { n = 4 }
		words := make([]string, n)
		for i := range words {
			idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(wordList))))
			words[i] = wordList[idx.Int64()]
		}
		fmt.Println(strings.Join(words, "-"))
	default:
		fmt.Fprintf(os.Stderr, "keygen: unknown type %q\n", typ); usage()
	}
}
