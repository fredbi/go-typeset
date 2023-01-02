package wordbreaker

type (
	// WordBreaker knows how to break words in parts.
	WordBreaker interface {
		BreakWord(string) []string
	}

	// SplitFunc splits a string into parts
	SplitFunc func(string) []string
)
