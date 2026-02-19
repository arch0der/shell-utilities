// escape - escape/unescape special characters in strings
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func escapeShell(s string) string {
	if s == "" { return "''" }
	needsQuote := false
	for _, ch := range s {
		if strings.ContainsRune(" \t\n\"'\\!$&;|<>(){}#~", ch) { needsQuote = true; break }
	}
	if !needsQuote { return s }
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func escapeRegex(s string) string {
	meta := `\.+*?()|[]{}^$`
	var b strings.Builder
	for _, ch := range s {
		if strings.ContainsRune(meta, ch) { b.WriteRune('\\') }
		b.WriteRune(ch)
	}
	return b.String()
}

func escapeSQL(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `'`, `''`)
}

func unescapeGo(s string) string {
	// Use strconv to interpret Go escape sequences
	if !strings.HasPrefix(s, `"`) { s = `"` + s + `"` }
	v, err := strconv.Unquote(s)
	if err != nil { return s }
	return v
}

func escapeGo(s string) string { return strconv.Quote(s) }

func escapeXML(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;",
	)
	return r.Replace(s)
}

func unescapeXML(s string) string {
	r := strings.NewReplacer(
		"&amp;", "&", "&lt;", "<", "&gt;", ">", "&quot;", `"`, "&#39;", "'",
	)
	return r.Replace(s)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: escape <mode> [text...]")
	fmt.Fprintln(os.Stderr, "  modes: shell | regex | sql | go | ungo | xml | unxml")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 { usage() }
	mode := os.Args[1]

	process := func(s string) string {
		switch mode {
		case "shell": return escapeShell(s)
		case "regex": return escapeRegex(s)
		case "sql": return escapeSQL(s)
		case "go": return escapeGo(s)
		case "ungo": return unescapeGo(s)
		case "xml": return escapeXML(s)
		case "unxml": return unescapeXML(s)
		default: fmt.Fprintf(os.Stderr, "escape: unknown mode %q\n", mode); os.Exit(1)
		}
		return s
	}

	if len(os.Args) > 2 {
		fmt.Println(process(strings.Join(os.Args[2:], " ")))
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() { fmt.Println(process(sc.Text())) }
}
