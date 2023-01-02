package hyphenator

import (
	"sync"
	"unicode"

	"golang.org/x/text/runes"
)

var (
	mx sync.Mutex

	// A cache for dictionaries. Users of this package only pay the cost
	// of building the trie once for a given language.
	loadedPatterns map[string]*Dictionary

	// hyphens is the set of unicode runes in range Hyphen
	hyphens = runes.In(unicode.Hyphen)
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

// load a preloaded trie Dictionary for some language patterns file.
func loadDictFromCache(patterns string) *Dictionary {
	mx.Lock()
	defer mx.Unlock()

	if loadedPatterns == nil {
		loadedPatterns = make(map[string]*Dictionary)
	}

	dict, ok := loadedPatterns[patterns]
	if !ok {
		dict, _ = LoadPatterns(patterns)
	}

	loadedPatterns[patterns] = dict

	return dict
}

// BreakWord breaks a word in parts which represent legitimate hyphenation points.
func (h *Hyphenator) BreakWord(word string) []string {
	if containsHyphen(word) {
		return []string{word} // TODO???
	}

	breakPoints, found := h.isException(word)
	if found {
		// known hyphenation rule exception
		return splitAtPositions(word, breakPoints)
	}

	result := make([]string, 0, 4)
	var positions = make([]int, 10) // the resulting hyphenation positions
	dottedWord := "." + word + "."
	for i := 0; i < len(dottedWord); i++ { // "word", "ord", "rd", "d"
		allBreakPoints, ok := h.isPattern(dottedWord[i:])
		if ok {
			positions = mergeBreakPoints(positions, allBreakPoints, i)
		}
	}

	positions = positions[1 : len(positions)-1]

	// ensure provided positions are consistent
	switch {
	case positions[0] > 0: // sometimes hyphen before first letter is "allowed"
		positions[0] = 0
	case len(positions) > len(word) && positions[len(word)] > 0:
		positions[len(word)] = 0 // sometimes hyphen after last letter is "allowed"
	}

	result = append(result, splitAtPositions(word, positions)...)

	return result
}

func (h *Hyphenator) isException(word string) ([]int, bool) {
	val := h.exceptions.Get(word)
	if val != nil {
		return val.([]int), true
	}

	return nil, false
}

func (h *Hyphenator) isPattern(fragment string) ([][]int, bool) {
	var result [][]int

	for j := 1; j < len(fragment); j++ {
		val := h.patterns.Get(fragment[:j])
		if val == nil {
			continue
		}
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

// split a string at given positions
func splitAtPositions(word string, positions []int) []string {
	parts := make([]string, 0, len(positions))
	previous := 0

	for i, pos := range positions {
		// odd numbers stand for possible break points, even numbers forbidden ones.
		if pos == 0 || pos%2 == 0 {
			continue
		}

		parts = append(parts, word[previous:i])
		previous = i
	}

	parts = append(parts, word[previous:])

	return parts
}

func containsHyphen(word string) bool {
	for _, r := range word {
		if hyphens.Contains(r) {
			return true
		}
	}

	return false
}
