package punctuator

import (
	"unicode"

	"golang.org/x/text/runes"
)

// punctuations is the set of unicode runes in range Punct
var hyphens = runes.In(unicode.Hyphen)

type Punctuator struct {
}

func New() *Punctuator {
	return &Punctuator{}
}

// BreakWord breaks words along punctuation separators such as ",", ".", ":", ";", "?", "!", "&"...
//
// It conforms to the unicode Punctuation class, with the addition of the "|" (pipe).
//
// Example:
//
//	"a;b,c.d:e_f\g-" => ["a/", "b/", "c|", "d-", "e_", "f\", "g-"]
func (p *Punctuator) BreakWord(word string) []string {
	return breakAtFunc(word, punctSplitFunc)
}

func punctSplitFunc(r rune) bool {
	return unicode.IsPunct(r) || r == '|'
}

// breakAtFunc works like strings.FieldsFunc, but retains separators.
//
// Break always happens _after_ the separator.
func breakAtFunc(word string, isBreak func(rune) bool) []string {
	parts := make([]string, 0, len(word))
	previous := 0

	for i, r := range word {
		if !isBreak(r) {
			continue
		}

		parts = append(parts, word[previous:i+1])
		previous = i + 1
	}

	if previous < len(word) {
		if len(word[previous:]) > 0 {
			parts = append(parts, word[previous:])
		}
	}

	return parts
}
