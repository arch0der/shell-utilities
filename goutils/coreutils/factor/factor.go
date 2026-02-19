package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func factorize(n *big.Int) []*big.Int {
	var factors []*big.Int
	two := big.NewInt(2)
	zero := big.NewInt(0)
	mod := new(big.Int)
	for new(big.Int).Mul(two, two).Cmp(n) <= 0 {
		for mod.Mod(n, two).Cmp(zero) == 0 {
			factors = append(factors, new(big.Int).Set(two))
			n.Div(n, two)
		}
		two.Add(two, big.NewInt(1))
	}
	if n.Cmp(big.NewInt(1)) > 0 {
		factors = append(factors, new(big.Int).Set(n))
	}
	return factors
}

func main() {
	args := os.Args[1:]
	process := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		n := new(big.Int)
		if _, ok := n.SetString(s, 10); !ok {
			fmt.Fprintf(os.Stderr, "factor: '%s' is not a valid positive integer\n", s)
			return
		}
		factors := factorize(new(big.Int).Set(n))
		parts := []string{n.String() + ":"}
		for _, f := range factors {
			parts = append(parts, f.String())
		}
		fmt.Println(strings.Join(parts, " "))
	}
	if len(args) == 0 {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			for _, tok := range strings.Fields(sc.Text()) {
				process(tok)
			}
		}
		return
	}
	for _, a := range args {
		process(a)
	}
}
