package tokenizer

import (
	"unicode"

	"github.com/fredbi/go-typeset/terminal/runes"
)

// Tokenizer breaks a text into space-separated tokens.
type Tokenizer struct {
}

// New tokenizer.
func New() *Tokenizer {
	return &Tokenizer{}
}

// BreakWord breaks a string into a slice of blank-separated tokens.
//
// Token separators are from the class unicode.IsSpace
//
// Blank separators are not retained in the result.
func (t *Tokenizer) BreakWord(word []rune) [][]rune {
	return runes.FieldsFunc(word, unicode.IsSpace)
}

// BreakWordString is the same as BreakWord but takes a string as input.
func (t *Tokenizer) BreakWordString(word string) [][]rune {
	return t.BreakWord([]rune(word))
}
