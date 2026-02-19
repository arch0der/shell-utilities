// bc - Basic calculator (arithmetic expression evaluator)
// Supports: +, -, *, /, %, ^, (, ), variables, sqrt(), sin(), cos(), etc.
// Usage: bc [-l] [file]  or pipe expressions to stdin
package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var mathLib = flag.Bool("l", false, "Load math library (sin, cos, etc.)")

var vars = map[string]float64{"pi": math.Pi, "e": math.E}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: bc [-l] [file]")
		flag.PrintDefaults()
	}
	flag.Parse()

	sc := bufio.NewScanner(os.Stdin)
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "bc:", err)
			os.Exit(1)
		}
		defer f.Close()
		sc = bufio.NewScanner(f)
	}

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || line == "quit" || line == "q" {
			if line == "quit" || line == "q" {
				break
			}
			continue
		}

		// Handle assignment: var = expr
		if idx := strings.Index(line, "="); idx > 0 {
			lhs := strings.TrimSpace(line[:idx])
			rhs := strings.TrimSpace(line[idx+1:])
			if isIdent(lhs) {
				val, err := eval(rhs)
				if err != nil {
					fmt.Fprintln(os.Stderr, "bc:", err)
					continue
				}
				vars[lhs] = val
				continue
			}
		}

		val, err := eval(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, "bc:", err)
			continue
		}
		if val == math.Trunc(val) && !strings.Contains(line, "/") {
			fmt.Printf("%d\n", int64(val))
		} else {
			fmt.Printf("%g\n", val)
		}
	}
}

func isIdent(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return len(s) > 0 && !unicode.IsDigit([]rune(s)[0])
}

// --- Tokenizer ---
type Token struct {
	kind int // 0=num, 1=ident, 2=op, 3=lparen, 4=rparen, 5=comma, 6=eof
	num  float64
	str  string
}

const (
	tNUM    = iota
	tIDENT
	tOP
	tLPAREN
	tRPAREN
	tCOMMA
	tEOF
)

func tokenize(expr string) []Token {
	var tokens []Token
	i := 0
	runes := []rune(expr)
	for i < len(runes) {
		r := runes[i]
		if unicode.IsSpace(r) {
			i++
			continue
		}
		if r == '#' {
			break
		}
		if r >= '0' && r <= '9' || r == '.' {
			start := i
			for i < len(runes) && (runes[i] >= '0' && runes[i] <= '9' || runes[i] == '.' || runes[i] == 'e' || runes[i] == 'E' || runes[i] == '+' && i > start || runes[i] == '-' && i > start) {
				i++
			}
			f, _ := strconv.ParseFloat(string(runes[start:i]), 64)
			tokens = append(tokens, Token{tNUM, f, ""})
			continue
		}
		if unicode.IsLetter(r) || r == '_' {
			start := i
			for i < len(runes) && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) || runes[i] == '_') {
				i++
			}
			tokens = append(tokens, Token{tIDENT, 0, string(runes[start:i])})
			continue
		}
		switch r {
		case '(':
			tokens = append(tokens, Token{tLPAREN, 0, "("})
		case ')':
			tokens = append(tokens, Token{tRPAREN, 0, ")"})
		case ',':
			tokens = append(tokens, Token{tCOMMA, 0, ","})
		default:
			// Check for two-char ops
			if i+1 < len(runes) {
				two := string(runes[i : i+2])
				if two == "**" {
					tokens = append(tokens, Token{tOP, 0, "**"})
					i += 2
					continue
				}
			}
			tokens = append(tokens, Token{tOP, 0, string(r)})
		}
		i++
	}
	tokens = append(tokens, Token{tEOF, 0, ""})
	return tokens
}

// --- Parser ---
type Parser struct {
	tokens []Token
	pos    int
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{tEOF, 0, ""}
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() Token {
	t := p.peek()
	p.pos++
	return t
}

func (p *Parser) parseExpr() (float64, error) {
	return p.parseAdd()
}

func (p *Parser) parseAdd() (float64, error) {
	left, err := p.parseMul()
	if err != nil {
		return 0, err
	}
	for {
		op := p.peek()
		if op.kind != tOP || (op.str != "+" && op.str != "-") {
			break
		}
		p.next()
		right, err := p.parseMul()
		if err != nil {
			return 0, err
		}
		if op.str == "+" {
			left += right
		} else {
			left -= right
		}
	}
	return left, nil
}

func (p *Parser) parseMul() (float64, error) {
	left, err := p.parsePow()
	if err != nil {
		return 0, err
	}
	for {
		op := p.peek()
		if op.kind != tOP || (op.str != "*" && op.str != "/" && op.str != "%" && op.str != "//") {
			break
		}
		p.next()
		right, err := p.parsePow()
		if err != nil {
			return 0, err
		}
		switch op.str {
		case "*":
			left *= right
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left /= right
		case "%":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left = math.Mod(left, right)
		case "//":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left = math.Floor(left / right)
		}
	}
	return left, nil
}

func (p *Parser) parsePow() (float64, error) {
	left, err := p.parseUnary()
	if err != nil {
		return 0, err
	}
	op := p.peek()
	if op.kind == tOP && (op.str == "^" || op.str == "**") {
		p.next()
		right, err := p.parsePow() // right-associative
		if err != nil {
			return 0, err
		}
		return math.Pow(left, right), nil
	}
	return left, nil
}

func (p *Parser) parseUnary() (float64, error) {
	op := p.peek()
	if op.kind == tOP && op.str == "-" {
		p.next()
		v, err := p.parsePrimary()
		return -v, err
	}
	if op.kind == tOP && op.str == "+" {
		p.next()
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (float64, error) {
	tok := p.peek()
	switch tok.kind {
	case tNUM:
		p.next()
		return tok.num, nil
	case tLPAREN:
		p.next()
		v, err := p.parseExpr()
		if err != nil {
			return 0, err
		}
		if p.peek().kind == tRPAREN {
			p.next()
		}
		return v, nil
	case tIDENT:
		name := tok.str
		p.next()
		if p.peek().kind == tLPAREN {
			p.next()
			var args []float64
			for p.peek().kind != tRPAREN && p.peek().kind != tEOF {
				arg, err := p.parseExpr()
				if err != nil {
					return 0, err
				}
				args = append(args, arg)
				if p.peek().kind == tCOMMA {
					p.next()
				}
			}
			if p.peek().kind == tRPAREN {
				p.next()
			}
			return callMath(name, args)
		}
		if v, ok := vars[name]; ok {
			return v, nil
		}
		return 0, fmt.Errorf("undefined variable: %s", name)
	}
	return 0, fmt.Errorf("unexpected token: %v", tok)
}

func callMath(name string, args []float64) (float64, error) {
	arg := func(i int) float64 {
		if i < len(args) {
			return args[i]
		}
		return 0
	}
	switch name {
	case "sqrt":
		return math.Sqrt(arg(0)), nil
	case "sin":
		return math.Sin(arg(0)), nil
	case "cos":
		return math.Cos(arg(0)), nil
	case "tan":
		return math.Tan(arg(0)), nil
	case "asin":
		return math.Asin(arg(0)), nil
	case "acos":
		return math.Acos(arg(0)), nil
	case "atan":
		return math.Atan(arg(0)), nil
	case "atan2":
		return math.Atan2(arg(0), arg(1)), nil
	case "log":
		return math.Log(arg(0)), nil
	case "log2":
		return math.Log2(arg(0)), nil
	case "log10":
		return math.Log10(arg(0)), nil
	case "exp":
		return math.Exp(arg(0)), nil
	case "pow":
		return math.Pow(arg(0), arg(1)), nil
	case "abs":
		return math.Abs(arg(0)), nil
	case "floor":
		return math.Floor(arg(0)), nil
	case "ceil":
		return math.Ceil(arg(0)), nil
	case "round":
		return math.Round(arg(0)), nil
	case "min":
		if len(args) < 2 {
			return 0, fmt.Errorf("min requires 2 args")
		}
		return math.Min(arg(0), arg(1)), nil
	case "max":
		if len(args) < 2 {
			return 0, fmt.Errorf("max requires 2 args")
		}
		return math.Max(arg(0), arg(1)), nil
	case "mod":
		return math.Mod(arg(0), arg(1)), nil
	case "gcd":
		a, b := uint64(math.Abs(arg(0))), uint64(math.Abs(arg(1)))
		for b != 0 {
			a, b = b, a%b
		}
		return float64(a), nil
	case "lcm":
		a, b := uint64(math.Abs(arg(0))), uint64(math.Abs(arg(1)))
		g := a
		bb := b
		for bb != 0 {
			g, bb = bb, g%bb
		}
		return float64(a * b / g), nil
	case "hypot":
		return math.Hypot(arg(0), arg(1)), nil
	case "factorial":
		n := int(arg(0))
		result := 1.0
		for i := 2; i <= n; i++ {
			result *= float64(i)
		}
		return result, nil
	}
	return 0, fmt.Errorf("unknown function: %s", name)
}

func eval(expr string) (float64, error) {
	tokens := tokenize(expr)
	p := &Parser{tokens: tokens}
	v, err := p.parseExpr()
	return v, err
}
