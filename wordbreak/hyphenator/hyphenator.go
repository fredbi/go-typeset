package hyphenator

import (
	"unicode"

	iface "github.com/fredbi/go-typeset/wordbreak"
	"golang.org/x/text/runes"
)

var (
	// hyphens is the set of unicode runes in range Hyphen
	hyphens = runes.In(unicode.Hyphen)

	_ iface.WordBreaker = &Hyphenator{}
)

// Hyphenator knows how to break words at legit hyphenation points.
type Hyphenator struct {
	*Dictionary
	*options
}

// New hyphenator.
func New(opts ...Option) *Hyphenator {
	h := &Hyphenator{
		options: defaultOptions(opts),
	}

	h.Dictionary = loadDictFromCache(langToPattern(h.lang))

	return h
}

// BreakWordString does the same as BreakWord but takes a string as input.
func (h *Hyphenator) BreakWordString(word string) [][]rune {
	return h.BreakWord([]rune(word))
}

// BreakWord breaks a word in parts which represent legitimate hyphenation points.
//
// Hyphenation marks are not rendered.
//
// Example (for US English):
//
//	"example" => ["ex", "am", "ple"]
func (h *Hyphenator) BreakWord(word []rune) [][]rune {
	wordLength := len(word)

	if containsHyphen(word) || wordLength < h.minLength {
		return [][]rune{word} // words with hyphens already set are not broken down // TODO
	}

	breakPoints, found := h.isException(word)
	if found {
		// known hyphenation rule exception
		return h.splitAtPositions(word, wordLength, breakPoints)
	}

	// allocate buffers once for all iterations
	hasPatterns := false
	allBreakPoints := make([][]int, 0, wordLength+2)
	positions := make([]int, 30) // the resulting hyphenation positions. A reasonable size is preallocated.
	fragmentBuffer := make([]rune, 0, wordLength+2)

	for i := 0; i < wordLength; i++ { // ".word.", "ord.", "rd.", "d."
		allBreakPoints, hasPatterns = h.isPattern(i, word, fragmentBuffer, allBreakPoints)
		if hasPatterns {
			positions = mergeBreakPoints(positions, allBreakPoints, i)
		}
	}

	positions = positions[1 : len(positions)-1]

	// ensure provided positions are consistent
	// ex: [0 1 2 0 0 0 3 0 0 1]
	switch {
	case positions[0] > 0: // sometimes hyphen before first letter is "allowed"
		positions[0] = 0
	case len(positions) > wordLength && positions[wordLength] > 0:
		positions[wordLength] = 0 // sometimes hyphen after last letter is "allowed"
	}

	return h.splitAtPositions(word, wordLength, positions)
}

func (h *Hyphenator) isException(word []rune) ([]int, bool) {
	val := h.exceptions.Get(word)
	if val != nil {
		return val.([]int), true
	}

	return nil, false
}

func (h *Hyphenator) isPattern(index int, word []rune, fragment []rune, result [][]int) ([][]int, bool) {
	// ".word." => ".w", ".wo", ".wor", ".word", ".word." ("." is skipped)
	// "word."  => "w", "wo", "wor", "word", "word."
	const dot = '.'

	// reset buffers, keep allocated memory
	var start int
	result = result[:0:cap(result)]
	fragment = fragment[:0:cap(fragment)]
	if index == 0 {
		fragment = append(fragment, dot)
		start = 0
	} else {
		start = index - 1
	}

	for i := start; i < len(word); i++ {
		r := unicode.ToLower(word[i])
		fragment = append(fragment, r)
		val := h.patterns.Get(fragment)
		if val == nil {
			continue
		}

		positions := val.([]int)
		result = append(result, positions)
	}

	fragment = append(fragment, dot)
	val := h.patterns.Get(fragment)
	if val != nil {
		positions := val.([]int)
		result = append(result, positions)
	}

	return result, len(result) > 0
}

// Merge a collection of positions arrays to a given positions array at
// a given index. Positions are overwritten, if a new position is greater
// than the old one. If the positions array isn't long enough, it will be
// enlarged.
//
// Example:
//
//	 with p = [0,2,0,0] and pp = { [1,7], [0,0,3] }
//
//	 after merge at position 1:
//		p = [0,2,7,3].
func mergeBreakPoints(positions []int, partialPositions [][]int, at int) []int {
	for _, partialPosition := range partialPositions {
		for relativeAt, num := range partialPosition { // for every relative position
			if missing := at + relativeAt - len(positions) + 1; missing > 0 {
				// grow positions
				for i := 0; i < missing; i++ {
					positions = append(positions, 0)
				}
			}

			if num > positions[at+relativeAt] { // new pos greater than current pos?
				positions[at+relativeAt] = num
			}
		}
	}

	return positions
}

// split a string at given positions: this applies the mask provided by positions to the cut the word.
func (h *Hyphenator) splitAtPositions(wordAsRunes []rune, length int, positions []int) [][]rune {
	parts := make([][]rune, 0, len(positions))
	previous := 0

	for i, pos := range positions {
		// odd numbers stand for possible break points, even numbers forbidden ones.
		if pos == 0 || pos%2 == 0 {
			continue
		}
		if i-previous < h.minLeft || length-i < h.minRight {
			continue
		}

		parts = append(parts, wordAsRunes[previous:i])
		previous = i
	}

	parts = append(parts, wordAsRunes[previous:])

	return parts
}

func containsHyphen(word []rune) bool {
	for _, r := range word {
		if hyphens.Contains(r) {
			return true
		}
	}

	return false
}

func IsHyphen(word []rune) bool {
	if len(word) != 1 {
		return false
	}

	return hyphens.Contains(word[0])
}

// SplitWord returns the parts of a word that are separated by an hyphen.
//
// Hyphens are isolated as individual parts. Using SplitWord in conjunction with BreakWord
// allows to process differently explicit hyphens from legit hyphens.
//
// Example:
//
//	"often-times" is transformed into ["often", "-", "times"].
func SplitWord(word []rune) [][]rune {
	result := make([][]rune, 0, len(word)+4)
	var previous int

	for i, r := range word {
		if hyphens.Contains(r) {
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
