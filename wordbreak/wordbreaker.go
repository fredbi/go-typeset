package wordbreaker

type (
	// WordBreaker knows how to break words in parts.
	WordBreaker interface {
		BreakWord([]rune) [][]rune
	}

	// SplitFunc splits a word into parts
	SplitFunc func([]rune) [][]rune
)
