// base58 - encode/decode Base58 and Base58Check (Bitcoin-style)
package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"strings"
)

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var bigZero = big.NewInt(0)
var bigBase = big.NewInt(58)

func encode(data []byte) string {
	n := new(big.Int).SetBytes(data)
	var result []byte
	mod := new(big.Int)
	for n.Cmp(bigZero) > 0 {
		n.DivMod(n, bigBase, mod)
		result = append(result, alphabet[mod.Int64()])
	}
	for _, b := range data { if b != 0 { break }; result = append(result, alphabet[0]) }
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 { result[i], result[j] = result[j], result[i] }
	return string(result)
}

func decode(s string) ([]byte, error) {
	n := big.NewInt(0)
	for _, ch := range s {
		idx := strings.IndexRune(alphabet, ch)
		if idx < 0 { return nil, fmt.Errorf("invalid character %q", ch) }
		n.Mul(n, bigBase)
		n.Add(n, big.NewInt(int64(idx)))
	}
	result := n.Bytes()
	for _, ch := range s { if ch != rune(alphabet[0]) { break }; result = append([]byte{0}, result...) }
	return result, nil
}

func checksum(data []byte) []byte {
	h1 := sha256.Sum256(data); h2 := sha256.Sum256(h1[:]); return h2[:4]
}

func encodeCheck(data []byte) string {
	payload := append(data, checksum(data)...); return encode(payload)
}

func decodeCheck(s string) ([]byte, error) {
	data, err := decode(s)
	if err != nil { return nil, err }
	if len(data) < 4 { return nil, fmt.Errorf("too short") }
	payload, cs := data[:len(data)-4], data[len(data)-4:]
	if string(checksum(payload)) != string(cs) { return nil, fmt.Errorf("checksum mismatch") }
	return payload, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: base58 <encode|decode|check-encode|check-decode> <hex_or_string>")
	os.Exit(1)
}

func hexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	if len(s)%2 != 0 { s = "0" + s }
	b := make([]byte, len(s)/2)
	for i := range b { fmt.Sscanf(s[2*i:2*i+2], "%02x", &b[i]) }
	return b
}

func main() {
	if len(os.Args) < 3 { usage() }
	cmd, arg := os.Args[1], os.Args[2]
	switch cmd {
	case "encode": fmt.Println(encode(hexToBytes(arg)))
	case "decode":
		b, err := decode(arg); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		fmt.Printf("%x\n", b)
	case "check-encode": fmt.Println(encodeCheck(hexToBytes(arg)))
	case "check-decode":
		b, err := decodeCheck(arg); if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
		fmt.Printf("%x\n", b)
	default: usage()
	}
}
