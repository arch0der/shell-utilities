// pidigits - print digits of pi, e, or golden ratio to arbitrary precision
package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// Machin's formula: pi/4 = 4*arctan(1/5) - arctan(1/239)
func arccot(x, prec uint) *big.Int {
	bigX := big.NewInt(int64(x))
	unity := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(prec+10)), nil)
	xSq := new(big.Int).Mul(bigX, bigX)
	sum := new(big.Int).Set(unity)
	sum.Div(sum, bigX)
	n := new(big.Int).Set(sum)
	sign := -1
	for i := 3; ; i += 2 {
		n.Div(n, xSq)
		if n.Sign() == 0 { break }
		term := new(big.Int).Div(n, big.NewInt(int64(i)))
		if sign < 0 { sum.Sub(sum, term) } else { sum.Add(sum, term) }
		sign = -sign
	}
	return sum
}

func computePi(digits uint) string {
	prec := digits + 10
	pi := new(big.Int)
	// pi = 4 * (4*arctan(1/5) - arctan(1/239))
	a := arccot(5, prec)
	b := arccot(239, prec)
	a.Mul(a, big.NewInt(4))
	a.Sub(a, b)
	a.Mul(a, big.NewInt(4))
	unity := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(prec+10)), nil)
	_ = unity
	pi.Set(a)
	s := pi.String()
	// Insert decimal point
	if len(s) > 1 { s = s[0:1] + "." + s[1:] }
	if uint(len(s)) > digits+2 { s = s[:digits+2] }
	return s
}

func computeE(digits int) string {
	prec := uint(digits + 20)
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(prec)), nil)
	result := new(big.Int).Set(scale) // 1
	term := new(big.Int).Set(scale)   // 1
	for i := int64(1); ; i++ {
		term.Div(term, big.NewInt(i))
		if term.Sign() == 0 { break }
		result.Add(result, term)
	}
	s := result.String()
	if len(s) > 1 { s = s[0:1] + "." + s[1:] }
	if len(s) > digits+2 { s = s[:digits+2] }
	return s
}

func main() {
	constant := "pi"
	digits := 50
	group := 10
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "pi","e","phi": constant = args[i]
		case "-g": i++; group, _ = strconv.Atoi(args[i])
		default:
			if n, err := strconv.Atoi(args[i]); err == nil { digits = n }
		}
	}
	if digits > 10000 { digits = 10000 }

	var result string
	switch constant {
	case "pi": result = computePi(uint(digits))
	case "e": result = computeE(digits)
	case "phi":
		// phi = (1 + sqrt(5)) / 2 — use big.Float
		prec := uint(digits * 4)
		sqrt5 := new(big.Float).SetPrec(prec).SetInt64(5)
		sqrt5.Sqrt(sqrt5)
		one := new(big.Float).SetPrec(prec).SetInt64(1)
		phi := new(big.Float).Add(one, sqrt5)
		phi.Quo(phi, new(big.Float).SetInt64(2))
		result = phi.Text('f', digits)
	}

	fmt.Printf("%s = ", constant)
	// Print with groups
	intPart := result[:strings.Index(result, ".")+1]
	fracPart := result[strings.Index(result, ".")+1:]
	fmt.Print(intPart)
	for i, ch := range fracPart {
		if i > 0 && i%group == 0 { fmt.Print(" ") }
		fmt.Printf("%c", ch)
	}
	fmt.Println()
}
