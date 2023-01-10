package tokenizer

import (
	"unicode"

	"github.com/fredbi/go-typeset/terminal/runes"
)

type Tokenizer struct {
}

func New() *Tokenizer {
	return &Tokenizer{}
}

// BreakWord breaks a string into a slice of blank-separated tokens.
//
// Blank separators are not retained in the result.
func (t *Tokenizer) BreakWord(word []rune) [][]rune {
	return runes.FieldsFunc(word, unicode.IsSpace)
}
