// awk - Pattern-action text processor (subset implementation)
// Supports: field splitting ($1..$NF, $0), NR, NF, FS, OFS,
//           BEGIN/END blocks, print, printf, if/else, arithmetic,
//           comparison operators, pattern matching (/regex/)
// Usage: awk [-F sep] [-v var=val] 'program' [file...]
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	fieldSep = flag.String("F", " ", "Field separator")
	vars     = flag.String("v", "", "Variable assignment (var=val)")
)

type Env struct {
	FS     string
	OFS    string
	ORS    string
	NR     int
	NF     int
	fields []string
	line   string
	vars   map[string]string
}

func newEnv(fs string) *Env {
	e := &Env{
		FS:   fs,
		OFS:  " ",
		ORS:  "\n",
		vars: map[string]string{},
	}
	return e
}

func (e *Env) setLine(line string) {
	e.line = line
	if e.FS == " " {
		e.fields = strings.Fields(line)
	} else {
		re := regexp.MustCompile(regexp.QuoteMeta(e.FS))
		e.fields = re.Split(line, -1)
	}
	e.NF = len(e.fields)
}

func (e *Env) getField(n int) string {
	if n == 0 {
		return e.line
	}
	if n >= 1 && n <= len(e.fields) {
		return e.fields[n-1]
	}
	return ""
}

func (e *Env) get(name string) string {
	switch name {
	case "NR":
		return strconv.Itoa(e.NR)
	case "NF":
		return strconv.Itoa(e.NF)
	case "FS":
		return e.FS
	case "OFS":
		return e.OFS
	case "ORS":
		return e.ORS
	}
	return e.vars[name]
}

func (e *Env) set(name, val string) {
	switch name {
	case "FS":
		e.FS = val
	case "OFS":
		e.OFS = val
	default:
		e.vars[name] = val
	}
}

// --- Tokenizer ---

type TokType int

const (
	TOK_NUM TokType = iota
	TOK_STR
	TOK_IDENT
	TOK_FIELD // $N
	TOK_OP
	TOK_LPAREN
	TOK_RPAREN
	TOK_LBRACE
	TOK_RBRACE
	TOK_COMMA
	TOK_SEMI
	TOK_REGEX
	TOK_EOF
)

type Token struct {
	typ TokType
	val string
}

type Tokenizer struct {
	src []rune
	pos int
}

func newTok(src string) *Tokenizer { return &Tokenizer{src: []rune(src)} }

func (t *Tokenizer) peek() rune {
	if t.pos >= len(t.src) {
		return 0
	}
	return t.src[t.pos]
}

func (t *Tokenizer) next() rune {
	if t.pos >= len(t.src) {
		return 0
	}
	r := t.src[t.pos]
	t.pos++
	return r
}

func (t *Tokenizer) skipWS() {
	for t.pos < len(t.src) && (t.src[t.pos] == ' ' || t.src[t.pos] == '\t' || t.src[t.pos] == '\n') {
		t.pos++
	}
}

func (t *Tokenizer) nextToken() Token {
	t.skipWS()
	if t.pos >= len(t.src) {
		return Token{TOK_EOF, ""}
	}
	c := t.peek()

	// Comment
	if c == '#' {
		for t.pos < len(t.src) && t.src[t.pos] != '\n' {
			t.pos++
		}
		return t.nextToken()
	}

	// Number
	if c >= '0' && c <= '9' || (c == '.' && t.pos+1 < len(t.src) && t.src[t.pos+1] >= '0' && t.src[t.pos+1] <= '9') {
		start := t.pos
		for t.pos < len(t.src) && (t.src[t.pos] >= '0' && t.src[t.pos] <= '9' || t.src[t.pos] == '.') {
			t.pos++
		}
		return Token{TOK_NUM, string(t.src[start:t.pos])}
	}

	// String
	if c == '"' {
		t.next()
		var sb strings.Builder
		for t.pos < len(t.src) && t.src[t.pos] != '"' {
			ch := t.next()
			if ch == '\\' && t.pos < len(t.src) {
				nxt := t.next()
				switch nxt {
				case 'n':
					sb.WriteByte('\n')
				case 't':
					sb.WriteByte('\t')
				default:
					sb.WriteRune(nxt)
				}
			} else {
				sb.WriteRune(ch)
			}
		}
		t.next() // closing "
		return Token{TOK_STR, sb.String()}
	}

	// Regex
	if c == '/' {
		t.next()
		var sb strings.Builder
		for t.pos < len(t.src) && t.src[t.pos] != '/' {
			sb.WriteRune(t.next())
		}
		t.next()
		return Token{TOK_REGEX, sb.String()}
	}

	// Field $N
	if c == '$' {
		t.next()
		start := t.pos
		for t.pos < len(t.src) && t.src[t.pos] >= '0' && t.src[t.pos] <= '9' {
			t.pos++
		}
		return Token{TOK_FIELD, string(t.src[start:t.pos])}
	}

	// Identifier / keyword
	if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' {
		start := t.pos
		for t.pos < len(t.src) && (t.src[t.pos] >= 'a' && t.src[t.pos] <= 'z' || t.src[t.pos] >= 'A' && t.src[t.pos] <= 'Z' || t.src[t.pos] >= '0' && t.src[t.pos] <= '9' || t.src[t.pos] == '_') {
			t.pos++
		}
		return Token{TOK_IDENT, string(t.src[start:t.pos])}
	}

	// Braces/parens
	switch c {
	case '(':
		t.next()
		return Token{TOK_LPAREN, "("}
	case ')':
		t.next()
		return Token{TOK_RPAREN, ")"}
	case '{':
		t.next()
		return Token{TOK_LBRACE, "{"}
	case '}':
		t.next()
		return Token{TOK_RBRACE, "}"}
	case ',':
		t.next()
		return Token{TOK_COMMA, ","}
	case ';':
		t.next()
		return Token{TOK_SEMI, ";"}
	}

	// Operators (multi-char)
	ops := []string{"==", "!=", "<=", ">=", "&&", "||", "++", "--", "+=", "-=", "*=", "/=", "!~", "~"}
	for _, op := range ops {
		if strings.HasPrefix(string(t.src[t.pos:]), op) {
			t.pos += len([]rune(op))
			return Token{TOK_OP, op}
		}
	}

	t.next()
	return Token{TOK_OP, string(c)}
}

func tokenize(src string) []Token {
	t := newTok(src)
	var toks []Token
	for {
		tok := t.nextToken()
		toks = append(toks, tok)
		if tok.typ == TOK_EOF {
			break
		}
	}
	return toks
}

// --- Parser / Evaluator ---

type Parser struct {
	toks []Token
	pos  int
	env  *Env
	out  io.Writer
}

func newParser(toks []Token, env *Env, out io.Writer) *Parser {
	return &Parser{toks: toks, env: env, out: out}
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.toks) {
		return Token{TOK_EOF, ""}
	}
	return p.toks[p.pos]
}

func (p *Parser) consume() Token {
	t := p.peek()
	p.pos++
	return t
}

func (p *Parser) expect(val string) {
	t := p.consume()
	if t.val != val {
		// silently continue
	}
}

// Parse the whole program: list of (pattern, action) rules
type Rule struct {
	pattern string // "BEGIN", "END", "/regex/", or expression string, or ""
	action  []Token
}

func parseProgram(src string) []Rule {
	toks := tokenize(src)
	var rules []Rule
	i := 0
	for i < len(toks) && toks[i].typ != TOK_EOF {
		// Skip semis
		if toks[i].typ == TOK_SEMI {
			i++
			continue
		}
		// Pattern
		pattern := ""
		if toks[i].typ != TOK_LBRACE {
			// Collect tokens until { 
			start := i
			depth := 0
			for i < len(toks) && !(toks[i].typ == TOK_LBRACE && depth == 0) {
				if toks[i].typ == TOK_LPAREN {
					depth++
				} else if toks[i].typ == TOK_RPAREN {
					depth--
				}
				i++
			}
			patToks := toks[start:i]
			parts := []string{}
			for _, t := range patToks {
				parts = append(parts, t.val)
			}
			pattern = strings.Join(parts, " ")
		}
		// Action
		if i < len(toks) && toks[i].typ == TOK_LBRACE {
			i++ // consume {
			start := i
			depth := 1
			for i < len(toks) && depth > 0 {
				if toks[i].typ == TOK_LBRACE {
					depth++
				} else if toks[i].typ == TOK_RBRACE {
					depth--
				}
				if depth > 0 {
					i++
				} else {
					break
				}
			}
			action := toks[start:i]
			i++ // consume }
			rules = append(rules, Rule{pattern: pattern, action: action})
		} else if pattern != "" {
			// Pattern without action = print
			rules = append(rules, Rule{pattern: pattern, action: tokenize("print $0")})
		}
	}
	return rules
}

func evalPattern(pattern string, env *Env) bool {
	if pattern == "" {
		return true
	}
	p := strings.TrimSpace(pattern)
	if p == "BEGIN" || p == "END" {
		return false
	}
	// Regex pattern /foo/
	if strings.HasPrefix(p, "/") && strings.HasSuffix(p, "/") {
		re := p[1 : len(p)-1]
		matched, _ := regexp.MatchString(re, env.line)
		return matched
	}
	// Expression
	toks := tokenize(p)
	ep := &ExprParser{toks: toks, env: env}
	val := ep.parseExpr()
	return isTruthy(val)
}

func isTruthy(s string) bool {
	if s == "" || s == "0" {
		return false
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f != 0
	}
	return true
}

func toFloat(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

func fmtNum(f float64) string {
	if f == math.Trunc(f) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'g', 6, 64)
}

// Expression parser
type ExprParser struct {
	toks []Token
	pos  int
	env  *Env
}

func (ep *ExprParser) peek() Token {
	if ep.pos >= len(ep.toks) {
		return Token{TOK_EOF, ""}
	}
	return ep.toks[ep.pos]
}

func (ep *ExprParser) consume() Token {
	t := ep.peek()
	ep.pos++
	return t
}

func (ep *ExprParser) parseExpr() string {
	return ep.parseOr()
}

func (ep *ExprParser) parseOr() string {
	left := ep.parseAnd()
	for ep.peek().val == "||" {
		ep.consume()
		right := ep.parseAnd()
		if isTruthy(left) || isTruthy(right) {
			left = "1"
		} else {
			left = "0"
		}
	}
	return left
}

func (ep *ExprParser) parseAnd() string {
	left := ep.parseCompare()
	for ep.peek().val == "&&" {
		ep.consume()
		right := ep.parseCompare()
		if isTruthy(left) && isTruthy(right) {
			left = "1"
		} else {
			left = "0"
		}
	}
	return left
}

func (ep *ExprParser) parseCompare() string {
	left := ep.parseConcat()
	for {
		op := ep.peek().val
		if op != "==" && op != "!=" && op != "<" && op != ">" && op != "<=" && op != ">=" && op != "~" && op != "!~" {
			break
		}
		ep.consume()
		right := ep.parseConcat()
		var result bool
		lf, lerr := strconv.ParseFloat(left, 64)
		rf, rerr := strconv.ParseFloat(right, 64)
		switch op {
		case "==":
			if lerr == nil && rerr == nil {
				result = lf == rf
			} else {
				result = left == right
			}
		case "!=":
			if lerr == nil && rerr == nil {
				result = lf != rf
			} else {
				result = left != right
			}
		case "<":
			if lerr == nil && rerr == nil {
				result = lf < rf
			} else {
				result = left < right
			}
		case ">":
			if lerr == nil && rerr == nil {
				result = lf > rf
			} else {
				result = left > right
			}
		case "<=":
			if lerr == nil && rerr == nil {
				result = lf <= rf
			} else {
				result = left <= right
			}
		case ">=":
			if lerr == nil && rerr == nil {
				result = lf >= rf
			} else {
				result = left >= right
			}
		case "~":
			matched, _ := regexp.MatchString(right, left)
			result = matched
		case "!~":
			matched, _ := regexp.MatchString(right, left)
			result = !matched
		}
		if result {
			left = "1"
		} else {
			left = "0"
		}
	}
	return left
}

func (ep *ExprParser) parseConcat() string {
	left := ep.parseAdd()
	// String concatenation: two adjacent values
	for {
		tok := ep.peek()
		if tok.typ == TOK_EOF || tok.typ == TOK_RBRACE || tok.typ == TOK_RPAREN ||
			tok.typ == TOK_SEMI || tok.typ == TOK_COMMA ||
			tok.val == ")" || tok.val == ";" || tok.val == "," ||
			tok.val == "}" || tok.val == "+" || tok.val == "-" ||
			tok.val == "*" || tok.val == "/" || tok.val == "%" ||
			tok.val == "==" || tok.val == "!=" || tok.val == "<" ||
			tok.val == ">" || tok.val == "<=" || tok.val == ">=" ||
			tok.val == "&&" || tok.val == "||" || tok.val == "=" ||
			tok.val == "+=" || tok.val == "-=" || tok.val == "~" || tok.val == "!~" {
			break
		}
		right := ep.parseAdd()
		left = left + right
	}
	return left
}

func (ep *ExprParser) parseAdd() string {
	left := ep.parseMul()
	for {
		op := ep.peek().val
		if op != "+" && op != "-" {
			break
		}
		ep.consume()
		right := ep.parseMul()
		lf, rf := toFloat(left), toFloat(right)
		if op == "+" {
			left = fmtNum(lf + rf)
		} else {
			left = fmtNum(lf - rf)
		}
	}
	return left
}

func (ep *ExprParser) parseMul() string {
	left := ep.parseUnary()
	for {
		op := ep.peek().val
		if op != "*" && op != "/" && op != "%" {
			break
		}
		ep.consume()
		right := ep.parseUnary()
		lf, rf := toFloat(left), toFloat(right)
		switch op {
		case "*":
			left = fmtNum(lf * rf)
		case "/":
			if rf == 0 {
				left = "inf"
			} else {
				left = fmtNum(lf / rf)
			}
		case "%":
			left = fmtNum(math.Mod(lf, rf))
		}
	}
	return left
}

func (ep *ExprParser) parseUnary() string {
	if ep.peek().val == "!" {
		ep.consume()
		v := ep.parsePrimary()
		if isTruthy(v) {
			return "0"
		}
		return "1"
	}
	if ep.peek().val == "-" {
		ep.consume()
		v := ep.parsePrimary()
		return fmtNum(-toFloat(v))
	}
	return ep.parsePrimary()
}

func (ep *ExprParser) parsePrimary() string {
	tok := ep.peek()
	switch tok.typ {
	case TOK_NUM:
		ep.consume()
		return tok.val
	case TOK_STR:
		ep.consume()
		return tok.val
	case TOK_FIELD:
		ep.consume()
		n, _ := strconv.Atoi(tok.val)
		return ep.env.getField(n)
	case TOK_REGEX:
		ep.consume()
		matched, _ := regexp.MatchString(tok.val, ep.env.line)
		if matched {
			return "1"
		}
		return "0"
	case TOK_LPAREN:
		ep.consume()
		v := ep.parseExpr()
		if ep.peek().typ == TOK_RPAREN {
			ep.consume()
		}
		return v
	case TOK_IDENT:
		ep.consume()
		name := tok.val
		// Function calls
		if ep.peek().typ == TOK_LPAREN {
			ep.consume()
			var args []string
			for ep.peek().typ != TOK_RPAREN && ep.peek().typ != TOK_EOF {
				args = append(args, ep.parseExpr())
				if ep.peek().val == "," {
					ep.consume()
				}
			}
			ep.consume() // )
			return callFunc(name, args, ep.env)
		}
		return ep.env.get(name)
	}
	return ""
}

func callFunc(name string, args []string, env *Env) string {
	switch name {
	case "length":
		if len(args) == 0 {
			return strconv.Itoa(len([]rune(env.line)))
		}
		return strconv.Itoa(len([]rune(args[0])))
	case "substr":
		if len(args) < 2 {
			return ""
		}
		s := []rune(args[0])
		start := int(toFloat(args[1])) - 1
		if start < 0 {
			start = 0
		}
		if start > len(s) {
			return ""
		}
		if len(args) >= 3 {
			n := int(toFloat(args[2]))
			end := start + n
			if end > len(s) {
				end = len(s)
			}
			return string(s[start:end])
		}
		return string(s[start:])
	case "index":
		if len(args) < 2 {
			return "0"
		}
		i := strings.Index(args[0], args[1])
		if i < 0 {
			return "0"
		}
		return strconv.Itoa(i + 1)
	case "split":
		if len(args) < 2 {
			return "0"
		}
		// split(s, arr, sep) - simplified
		return strconv.Itoa(len(strings.Fields(args[0])))
	case "sub", "gsub":
		if len(args) < 2 {
			return ""
		}
		re := regexp.MustCompile(args[0])
		if name == "gsub" {
			return re.ReplaceAllString(env.line, args[1])
		}
		return re.ReplaceAllLiteralString(env.line, args[1])
	case "match":
		if len(args) < 2 {
			return "0"
		}
		re := regexp.MustCompile(args[1])
		loc := re.FindStringIndex(args[0])
		if loc == nil {
			return "0"
		}
		return strconv.Itoa(loc[0] + 1)
	case "toupper":
		if len(args) < 1 {
			return ""
		}
		return strings.ToUpper(args[0])
	case "tolower":
		if len(args) < 1 {
			return ""
		}
		return strings.ToLower(args[0])
	case "int":
		if len(args) < 1 {
			return "0"
		}
		return strconv.Itoa(int(toFloat(args[0])))
	case "sqrt":
		if len(args) < 1 {
			return "0"
		}
		return fmtNum(math.Sqrt(toFloat(args[0])))
	case "sin":
		if len(args) < 1 {
			return "0"
		}
		return fmtNum(math.Sin(toFloat(args[0])))
	case "cos":
		if len(args) < 1 {
			return "0"
		}
		return fmtNum(math.Cos(toFloat(args[0])))
	case "exp":
		if len(args) < 1 {
			return "0"
		}
		return fmtNum(math.Exp(toFloat(args[0])))
	case "log":
		if len(args) < 1 {
			return "0"
		}
		return fmtNum(math.Log(toFloat(args[0])))
	case "atan2":
		if len(args) < 2 {
			return "0"
		}
		return fmtNum(math.Atan2(toFloat(args[0]), toFloat(args[1])))
	}
	return ""
}

// Statement executor
type StmtExec struct {
	toks []Token
	pos  int
	env  *Env
	out  io.Writer
}

func (se *StmtExec) peek() Token {
	if se.pos >= len(se.toks) {
		return Token{TOK_EOF, ""}
	}
	return se.toks[se.pos]
}

func (se *StmtExec) consume() Token {
	t := se.peek()
	se.pos++
	return t
}

func (se *StmtExec) evalExpr() string {
	// Find extent of expression (up to ; or } or newline context)
	start := se.pos
	depth := 0
	for se.pos < len(se.toks) {
		t := se.toks[se.pos]
		if t.typ == TOK_LPAREN {
			depth++
		} else if t.typ == TOK_RPAREN {
			if depth == 0 {
				break
			}
			depth--
		}
		if depth == 0 {
			if t.typ == TOK_SEMI || t.typ == TOK_RBRACE || t.typ == TOK_EOF {
				break
			}
			if t.val == "," {
				break
			}
		}
		se.pos++
	}
	ep := &ExprParser{toks: se.toks[start:se.pos], env: se.env}
	return ep.parseExpr()
}

func (se *StmtExec) run() {
	for se.pos < len(se.toks) && se.toks[se.pos].typ != TOK_EOF {
		se.runStmt()
	}
}

func (se *StmtExec) runStmt() {
	tok := se.peek()
	if tok.typ == TOK_SEMI {
		se.consume()
		return
	}
	if tok.typ == TOK_RBRACE || tok.typ == TOK_EOF {
		return
	}

	switch tok.val {
	case "print":
		se.consume()
		var parts []string
		for se.peek().typ != TOK_SEMI && se.peek().typ != TOK_RBRACE && se.peek().typ != TOK_EOF {
			parts = append(parts, se.evalExpr())
			if se.peek().val == "," {
				se.consume()
			}
		}
		if len(parts) == 0 {
			fmt.Fprintln(se.out, se.env.line)
		} else {
			fmt.Fprintln(se.out, strings.Join(parts, se.env.OFS))
		}

	case "printf":
		se.consume()
		var args []string
		for se.peek().typ != TOK_SEMI && se.peek().typ != TOK_RBRACE && se.peek().typ != TOK_EOF {
			args = append(args, se.evalExpr())
			if se.peek().val == "," {
				se.consume()
			}
		}
		if len(args) > 0 {
			fmtStr := args[0]
			fmtArgs := make([]interface{}, len(args)-1)
			for i, a := range args[1:] {
				if f, err := strconv.ParseFloat(a, 64); err == nil {
					fmtArgs[i] = f
				} else {
					fmtArgs[i] = a
				}
			}
			fmt.Fprintf(se.out, fmtStr, fmtArgs...)
		}

	case "if":
		se.consume()
		se.consume() // (
		cond := se.evalExpr()
		se.consume() // )
		se.consume() // {
		bodyToks := se.collectBlock()
		var elseToks []Token
		if se.peek().val == "else" {
			se.consume()
			if se.peek().typ == TOK_LBRACE {
				se.consume()
				elseToks = se.collectBlock()
			}
		}
		if isTruthy(cond) {
			sub := &StmtExec{toks: bodyToks, env: se.env, out: se.out}
			sub.run()
		} else if len(elseToks) > 0 {
			sub := &StmtExec{toks: elseToks, env: se.env, out: se.out}
			sub.run()
		}

	case "while":
		se.consume()
		se.consume() // (
		condStart := se.pos
		se.evalExpr()
		condEnd := se.pos
		se.consume() // )
		se.consume() // {
		bodyToks := se.collectBlock()
		for {
			ep := &ExprParser{toks: se.toks[condStart:condEnd], env: se.env}
			if !isTruthy(ep.parseExpr()) {
				break
			}
			sub := &StmtExec{toks: bodyToks, env: se.env, out: se.out}
			sub.run()
		}

	case "for":
		se.consume()
		se.consume() // (
		// init
		se.runStmt()
		if se.peek().typ == TOK_SEMI {
			se.consume()
		}
		condStart := se.pos
		se.evalExpr()
		condEnd := se.pos
		if se.peek().typ == TOK_SEMI {
			se.consume()
		}
		incrStart := se.pos
		se.evalExpr()
		incrEnd := se.pos
		se.consume() // )
		se.consume() // {
		bodyToks := se.collectBlock()
		for {
			ep := &ExprParser{toks: se.toks[condStart:condEnd], env: se.env}
			if !isTruthy(ep.parseExpr()) {
				break
			}
			sub := &StmtExec{toks: bodyToks, env: se.env, out: se.out}
			sub.run()
			incr := &StmtExec{toks: se.toks[incrStart:incrEnd], env: se.env, out: se.out}
			incr.runStmt()
		}

	case "next":
		se.consume()
		panic("next") // caught by caller

	default:
		// Assignment or expression
		if tok.typ == TOK_IDENT && se.pos+1 < len(se.toks) {
			next := se.toks[se.pos+1]
			if next.val == "=" || next.val == "+=" || next.val == "-=" || next.val == "*=" || next.val == "/=" {
				name := tok.val
				se.consume()
				op := se.consume().val
				val := se.evalExpr()
				cur := se.env.get(name)
				switch op {
				case "=":
					se.env.set(name, val)
				case "+=":
					se.env.set(name, fmtNum(toFloat(cur)+toFloat(val)))
				case "-=":
					se.env.set(name, fmtNum(toFloat(cur)-toFloat(val)))
				case "*=":
					se.env.set(name, fmtNum(toFloat(cur)*toFloat(val)))
				case "/=":
					se.env.set(name, fmtNum(toFloat(cur)/toFloat(val)))
				}
				return
			}
			if next.val == "++" {
				name := tok.val
				se.consume()
				se.consume()
				cur := toFloat(se.env.get(name))
				se.env.set(name, fmtNum(cur+1))
				return
			}
			if next.val == "--" {
				name := tok.val
				se.consume()
				se.consume()
				cur := toFloat(se.env.get(name))
				se.env.set(name, fmtNum(cur-1))
				return
			}
		}
		se.evalExpr()
	}

	if se.peek().typ == TOK_SEMI {
		se.consume()
	}
}

func (se *StmtExec) collectBlock() []Token {
	var toks []Token
	depth := 1
	for se.pos < len(se.toks) && depth > 0 {
		t := se.toks[se.pos]
		if t.typ == TOK_LBRACE {
			depth++
		} else if t.typ == TOK_RBRACE {
			depth--
			if depth == 0 {
				se.pos++
				break
			}
		}
		toks = append(toks, t)
		se.pos++
	}
	return toks
}

func execAction(toks []Token, env *Env, out io.Writer) (recovered bool) {
	defer func() {
		if r := recover(); r != nil {
			if r == "next" {
				recovered = true
			}
		}
	}()
	se := &StmtExec{toks: toks, env: env, out: out}
	se.run()
	return false
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: awk [-F sep] [-v var=val] 'program' [file...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	progSrc := flag.Arg(0)
	files := flag.Args()[1:]

	env := newEnv(*fieldSep)

	// Parse -v assignment
	if *vars != "" {
		parts := strings.SplitN(*vars, "=", 2)
		if len(parts) == 2 {
			env.set(parts[0], parts[1])
		}
	}

	rules := parseProgram(progSrc)

	// Execute BEGIN
	for _, rule := range rules {
		if rule.pattern == "BEGIN" {
			execAction(rule.action, env, os.Stdout)
		}
	}

	// Process files or stdin
	readers := []io.Reader{}
	if len(files) == 0 {
		readers = append(readers, os.Stdin)
	} else {
		for _, path := range files {
			f, err := os.Open(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "awk:", err)
				continue
			}
			defer f.Close()
			readers = append(readers, f)
		}
	}

	for _, r := range readers {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			env.NR++
			env.setLine(scanner.Text())
			for _, rule := range rules {
				if rule.pattern == "BEGIN" || rule.pattern == "END" {
					continue
				}
				if evalPattern(rule.pattern, env) {
					next := execAction(rule.action, env, os.Stdout)
					if next {
						break
					}
				}
			}
		}
	}

	// Execute END
	for _, rule := range rules {
		if rule.pattern == "END" {
			execAction(rule.action, env, os.Stdout)
		}
	}
}
