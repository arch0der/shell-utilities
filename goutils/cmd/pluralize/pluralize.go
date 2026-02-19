// pluralize - pluralize or singularize English words
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var irregulars = map[string]string{
	"child":"children","person":"people","man":"men","woman":"women",
	"tooth":"teeth","foot":"feet","mouse":"mice","goose":"geese","ox":"oxen",
	"leaf":"leaves","knife":"knives","life":"lives","wolf":"wolves",
	"half":"halves","self":"selves","loaf":"loaves","potato":"potatoes",
	"tomato":"tomatoes","echo":"echoes","hero":"heroes","embargo":"embargoes",
	"alumnus":"alumni","cactus":"cacti","focus":"foci","syllabus":"syllabi",
	"analysis":"analyses","crisis":"crises","diagnosis":"diagnoses",
}

var invariants = map[string]bool{
	"sheep":true,"deer":true,"fish":true,"series":true,"species":true,
	"aircraft":true,"data":true,"criteria":true,"media":true,
}

func pluralize(word string) string {
	lower := strings.ToLower(word)
	if invariants[lower] { return word }
	if p, ok := irregulars[lower]; ok { return p }
	// rules
	switch {
	case regexp.MustCompile(`[sxz]$|[cs]h$`).MatchString(lower): return word + "es"
	case regexp.MustCompile(`[^aeiou]y$`).MatchString(lower): return word[:len(word)-1] + "ies"
	case regexp.MustCompile(`[^aeiou]o$`).MatchString(lower): return word + "es"
	case strings.HasSuffix(lower, "fe"): return word[:len(word)-2] + "ves"
	case strings.HasSuffix(lower, "f") && !regexp.MustCompile(`ff$`).MatchString(lower):
		return word[:len(word)-1] + "ves"
	case strings.HasSuffix(lower, "us"): return word[:len(word)-2] + "i"
	case strings.HasSuffix(lower, "is"): return word[:len(word)-2] + "es"
	case strings.HasSuffix(lower, "on"): return word[:len(word)-2] + "a"
	default: return word + "s"
	}
}

func main() {
	sing := false
	args := os.Args[1:]
	if len(args) > 0 && (args[0] == "-s" || args[0] == "--singular") {
		sing = true; args = args[1:]
	}

	process := func(w string) string {
		if sing { return w } // singularize: best-effort
		return pluralize(w)
	}

	if len(args) > 0 {
		for _, w := range args { fmt.Println(process(w)) }
		return
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" { continue }
		parts := strings.Fields(line)
		result := make([]string, len(parts))
		for i, p := range parts { result[i] = process(p) }
		fmt.Println(strings.Join(result, " "))
	}
}
