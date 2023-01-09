package punctuator

import (
	"unicode"

	iface "github.com/fredbi/go-typeset/wordbreak"
	"golang.org/x/text/runes"
)

var (
	_       iface.WordBreaker = &Punctuator{}
	hyphens                   = runes.In(unicode.Hyphen)
)

// Punctuator knows how to break words on punctuation runes.
type Punctuator struct {
}

// New word breaker across punctuation marks.
func New() *Punctuator {
	return &Punctuator{}
}

// BreakWordString is like BreakWord but takes a string as input.
func (p *Punctuator) BreakWordString(word string) [][]rune {
	return p.BreakWord([]rune(word))
}

// BreakWord breaks words along punctuation separators such as ",", ".", ":", ";", "?", "!", "&"...
//
// Punctuation runes are retained as single character parts.
//
// It conforms to the unicode Punctuation class, with the addition of the "|" (pipe), but not hyphens (which are handled by the hyphenator package).
//
//	"a;b,c.d:e_f\g_" => ["a", ";", "b", ",", "c", ".", "d", ":", "e", "_", "f", "\", "g" ,"_"]
func (p *Punctuator) BreakWord(word []rune) [][]rune {
	return breakAtFunc(word, punctSplitFunc)
}

func punctSplitFunc(r rune) bool {
	return r == '|' || (unicode.IsPunct(r) && !hyphens.Contains(r))
}

// breakAtFunc works like strings.FieldsFunc, but retains separators.
func breakAtFunc(word []rune, isBreak func(rune) bool) [][]rune {
	result := make([][]rune, 0, len(word)+4)
	var previous int

	for i, r := range word {
		if isBreak(r) {
			if i > 0 && previous < len(word) && previous < i {
				result = append(result, word[previous:i])
			}
			result = append(result, []rune{r})
			previous = i + 1
		}
	}

	if previous < len(word) {
		result = append(result, word[previous:])
	}

	return result
}

func IsPunctuation(word []rune) bool {
	if len(word) != 1 {
		return false
	}

	return punctSplitFunc(word[0])
}
