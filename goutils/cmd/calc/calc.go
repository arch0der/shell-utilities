// calc - evaluate mathematical expressions from the command line
package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
)

// Recursive descent expression evaluator
type parser struct {
	s   scanner.Scanner
	tok rune
	lit string
}

func newParser(expr string) *parser {
	p := &parser{}
	p.s.Init(strings.NewReader(expr))
	p.s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanInts
	p.s.IsIdentRune = func(ch rune, i int) bool {
		return unicode.IsLetter(ch) || ch == '_' || (i > 0 && unicode.IsDigit(ch))
	}
	p.next()
	return p
}

func (p *parser) next() {
	p.tok = p.s.Scan()
	p.lit = p.s.TokenText()
}

var funcs = map[string]func(float64) float64{
	"sin": math.Sin, "cos": math.Cos, "tan": math.Tan,
	"asin": math.Asin, "acos": math.Acos, "atan": math.Atan,
	"sqrt": math.Sqrt, "cbrt": math.Cbrt, "abs": math.Abs,
	"log": math.Log10, "ln": math.Log, "log2": math.Log2,
	"exp": math.Exp, "ceil": math.Ceil, "floor": math.Floor,
	"round": math.Round,
}

var consts = map[string]float64{
	"pi": math.Pi, "e": math.E, "phi": math.Phi,
	"inf": math.Inf(1), "nan": math.NaN(),
}

func (p *parser) expr() float64 { return p.addSub() }

func (p *parser) addSub() float64 {
	x := p.mulDiv()
	for p.lit == "+" || p.lit == "-" {
		op := p.lit; p.next()
		y := p.mulDiv()
		if op == "+" { x += y } else { x -= y }
	}
	return x
}

func (p *parser) mulDiv() float64 {
	x := p.power()
	for p.lit == "*" || p.lit == "/" || p.lit == "%" {
		op := p.lit; p.next()
		y := p.power()
		switch op {
		case "*": x *= y
		case "/": if y == 0 { fmt.Fprintln(os.Stderr, "calc: division by zero"); os.Exit(1) }; x /= y
		case "%": x = math.Mod(x, y)
		}
	}
	return x
}

func (p *parser) power() float64 {
	x := p.unary()
	if p.lit == "^" || p.lit == "**" { p.next(); return math.Pow(x, p.power()) }
	return x
}

func (p *parser) unary() float64 {
	if p.lit == "-" { p.next(); return -p.primary() }
	if p.lit == "+" { p.next() }
	return p.primary()
}

func (p *parser) primary() float64 {
	if p.tok == scanner.Float || p.tok == scanner.Int {
		v, _ := strconv.ParseFloat(p.lit, 64); p.next(); return v
	}
	if p.tok == scanner.Ident {
		name := p.lit; p.next()
		if v, ok := consts[name]; ok { return v }
		if fn, ok := funcs[name]; ok {
			if p.lit != "(" { fmt.Fprintf(os.Stderr, "calc: expected ( after %s\n", name); os.Exit(1) }
			p.next(); v := fn(p.expr())
			if p.lit != ")" { fmt.Fprintln(os.Stderr, "calc: expected )"); os.Exit(1) }
			p.next(); return v
		}
		fmt.Fprintf(os.Stderr, "calc: unknown identifier %q\n", name); os.Exit(1)
	}
	if p.lit == "(" { p.next(); v := p.expr()
		if p.lit != ")" { fmt.Fprintln(os.Stderr, "calc: expected )"); os.Exit(1) }
		p.next(); return v
	}
	fmt.Fprintf(os.Stderr, "calc: unexpected %q\n", p.lit); os.Exit(1)
	return 0
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: calc <expression>")
		fmt.Fprintln(os.Stderr, "  example: calc '2^10 + sin(pi/2) * 100'")
		fmt.Fprintln(os.Stderr, "  consts: pi e phi")
		fmt.Fprintln(os.Stderr, "  funcs: sin cos tan asin acos atan sqrt log ln exp abs ceil floor round")
		os.Exit(1)
	}
	expr := strings.Join(os.Args[1:], "")
	p := newParser(expr)
	result := p.expr()
	if result == math.Trunc(result) && !math.IsInf(result, 0) {
		fmt.Printf("%.0f\n", result)
	} else {
		fmt.Printf("%g\n", result)
	}
}
