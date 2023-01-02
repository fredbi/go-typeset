package separator

import (
	"unicode"

	"golang.org/x/text/runes"
)

// hyphens is the set of unicode runes in range Hyphen
var hyphens = runes.In(unicode.Hyphen)

type Separator struct {
}

func New() *Separator {
	return &Separator{}
}

// BreakWord breaks words along "natural" separators such as "-", "_", "|", "/"...
//
// Example:
//
//	"a/b/c|d-e_f\g-" => ["a/", "b/", "c|", "d-", "e_", "f\", "g-"]
func (h *Separator) BreakWord(word string) []string {
	return breakAtFunc(word, separatorSplitFunc)
}

func separatorSplitFunc(r rune) bool {
	if hyphens.Contains(r) {
		return true
	}

	switch r {
	case '|', '/', '\\', '_':
		return true
	default:
		return false
	}
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
