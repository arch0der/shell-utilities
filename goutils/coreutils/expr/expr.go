package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "expr: missing operand")
		os.Exit(1)
	}
	result, _ := exprEval(args)
	fmt.Println(result)
	if result == "0" || result == "" {
		os.Exit(1)
	}
}

func exprEval(tokens []string) (string, []string) {
	return exprOr(tokens)
}

func exprOr(tokens []string) (string, []string) {
	left, tokens := exprAnd(tokens)
	for len(tokens) > 0 && tokens[0] == "|" {
		tokens = tokens[1:]
		right, rest := exprAnd(tokens)
		tokens = rest
		if left == "0" || left == "" {
			left = right
		}
	}
	return left, tokens
}

func exprAnd(tokens []string) (string, []string) {
	left, tokens := exprCmp(tokens)
	for len(tokens) > 0 && tokens[0] == "&" {
		tokens = tokens[1:]
		right, rest := exprCmp(tokens)
		tokens = rest
		if left == "0" || left == "" || right == "0" || right == "" {
			left = "0"
		}
	}
	return left, tokens
}

func exprCmp(tokens []string) (string, []string) {
	left, tokens := exprAdd(tokens)
	ops := []string{"<=", ">=", "!=", "=", "<", ">"}
	for len(tokens) > 0 {
		op := ""
		for _, o := range ops {
			if tokens[0] == o {
				op = o
				break
			}
		}
		if op == "" {
			break
		}
		tokens = tokens[1:]
		right, rest := exprAdd(tokens)
		tokens = rest
		ln, lerr := strconv.ParseInt(left, 10, 64)
		rn, rerr := strconv.ParseInt(right, 10, 64)
		var result bool
		if lerr == nil && rerr == nil {
			switch op {
			case "<":
				result = ln < rn
			case "<=":
				result = ln <= rn
			case ">":
				result = ln > rn
			case ">=":
				result = ln >= rn
			case "=":
				result = ln == rn
			case "!=":
				result = ln != rn
			}
		} else {
			switch op {
			case "<":
				result = left < right
			case "<=":
				result = left <= right
			case ">":
				result = left > right
			case ">=":
				result = left >= right
			case "=":
				result = left == right
			case "!=":
				result = left != right
			}
		}
		if result {
			left = "1"
		} else {
			left = "0"
		}
	}
	return left, tokens
}

func exprAdd(tokens []string) (string, []string) {
	left, tokens := exprMul(tokens)
	for len(tokens) > 0 && (tokens[0] == "+" || tokens[0] == "-") {
		op := tokens[0]
		tokens = tokens[1:]
		right, rest := exprMul(tokens)
		tokens = rest
		ln, _ := strconv.ParseInt(left, 10, 64)
		rn, _ := strconv.ParseInt(right, 10, 64)
		if op == "+" {
			left = strconv.FormatInt(ln+rn, 10)
		} else {
			left = strconv.FormatInt(ln-rn, 10)
		}
	}
	return left, tokens
}

func exprMul(tokens []string) (string, []string) {
	left, tokens := exprUnary(tokens)
	for len(tokens) > 0 && (tokens[0] == "*" || tokens[0] == "/" || tokens[0] == "%") {
		op := tokens[0]
		tokens = tokens[1:]
		right, rest := exprUnary(tokens)
		tokens = rest
		ln, _ := strconv.ParseInt(left, 10, 64)
		rn, _ := strconv.ParseInt(right, 10, 64)
		switch op {
		case "*":
			left = strconv.FormatInt(ln*rn, 10)
		case "/":
			if rn == 0 {
				fmt.Fprintln(os.Stderr, "expr: division by zero")
				os.Exit(2)
			}
			left = strconv.FormatInt(ln/rn, 10)
		case "%":
			if rn == 0 {
				fmt.Fprintln(os.Stderr, "expr: division by zero")
				os.Exit(2)
			}
			left = strconv.FormatInt(ln%rn, 10)
		}
	}
	return left, tokens
}

func exprUnary(tokens []string) (string, []string) {
	if len(tokens) == 0 {
		return "", tokens
	}
	if tokens[0] == "match" {
		tokens = tokens[1:]
		s, tokens := exprUnary(tokens)
		pat, tokens := exprUnary(tokens)
		re := regexp.MustCompile("^(?:" + pat + ")")
		m := re.FindStringSubmatch(s)
		if m == nil {
			return "0", tokens
		}
		if len(m) > 1 {
			return m[1], tokens
		}
		return strconv.Itoa(len(m[0])), tokens
	}
	if tokens[0] == "length" {
		tokens = tokens[1:]
		s, tokens := exprUnary(tokens)
		return strconv.Itoa(len(s)), tokens
	}
	if tokens[0] == "substr" {
		tokens = tokens[1:]
		s, tokens := exprUnary(tokens)
		posStr, tokens := exprUnary(tokens)
		lenStr, tokens := exprUnary(tokens)
		pos, _ := strconv.Atoi(posStr)
		length, _ := strconv.Atoi(lenStr)
		if pos < 1 {
			pos = 1
		}
		pos--
		runes := []rune(s)
		if pos >= len(runes) {
			return "", tokens
		}
		end := pos + length
		if end > len(runes) {
			end = len(runes)
		}
		return string(runes[pos:end]), tokens
	}
	if tokens[0] == "index" {
		tokens = tokens[1:]
		s, tokens := exprUnary(tokens)
		chars, tokens := exprUnary(tokens)
		idx := strings.IndexAny(s, chars)
		return strconv.Itoa(idx + 1), tokens
	}
	if tokens[0] == ":" {
		// match operator
		return "", tokens[1:]
	}
	val := tokens[0]
	return val, tokens[1:]
}
