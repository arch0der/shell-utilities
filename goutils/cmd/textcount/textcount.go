// textcount - count sentences, paragraphs, avg word length, reading time
package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var sentenceEnd = regexp.MustCompile(`[.!?]+[\s"')\]]*(\s|$)`)

func analyze(text string) {
	// Words
	words := strings.FieldsFunc(text, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	wordCount := len(words)

	// Characters
	charCount := len([]rune(text))
	charNoSpace := 0
	for _, r := range text { if !unicode.IsSpace(r) { charNoSpace++ } }

	// Sentences
	sentences := sentenceEnd.FindAllString(text, -1)
	sentenceCount := len(sentences)
	if sentenceCount == 0 && wordCount > 0 { sentenceCount = 1 }

	// Paragraphs
	paragraphs := 0
	for _, para := range strings.Split(text, "\n\n") {
		if strings.TrimSpace(para) != "" { paragraphs++ }
	}
	if paragraphs == 0 { paragraphs = 1 }

	// Avg word length
	totalLen := 0
	for _, w := range words { totalLen += len([]rune(w)) }
	avgWordLen := 0.0
	if wordCount > 0 { avgWordLen = float64(totalLen) / float64(wordCount) }

	// Flesch-Kincaid readability (approx)
	syllables := 0
	for _, w := range words { syllables += countSyllables(w) }
	var fkGrade float64
	if wordCount > 0 && sentenceCount > 0 {
		fkGrade = 0.39*(float64(wordCount)/float64(sentenceCount)) + 11.8*(float64(syllables)/float64(wordCount)) - 15.59
		fkGrade = math.Round(fkGrade*10) / 10
	}

	// Reading time (avg 238 wpm)
	mins := float64(wordCount) / 238.0
	readTime := fmt.Sprintf("%.0f min", math.Ceil(mins))
	if mins < 1 { readTime = fmt.Sprintf("%.0f sec", math.Ceil(mins*60)) }

	fmt.Printf("Words          : %d\n", wordCount)
	fmt.Printf("Characters     : %d  (%d without spaces)\n", charCount, charNoSpace)
	fmt.Printf("Sentences      : %d\n", sentenceCount)
	fmt.Printf("Paragraphs     : %d\n", paragraphs)
	fmt.Printf("Syllables      : %d\n", syllables)
	fmt.Printf("Avg word length: %.1f chars\n", avgWordLen)
	fmt.Printf("FK Grade Level : %.1f\n", fkGrade)
	fmt.Printf("Reading time   : ~%s\n", readTime)
}

func countSyllables(word string) int {
	word = strings.ToLower(word)
	count := 0
	prev := false
	for _, r := range word {
		vowel := strings.ContainsRune("aeiouy", r)
		if vowel && !prev { count++ }
		prev = vowel
	}
	if strings.HasSuffix(word, "e") && count > 1 { count-- }
	if count == 0 { count = 1 }
	return count
}

func main() {
	var data []byte
	var err error
	if len(os.Args) > 1 {
		data, err = os.ReadFile(os.Args[1])
	} else {
		data, err = io.ReadAll(os.Stdin)
	}
	if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
	analyze(string(data))
}
