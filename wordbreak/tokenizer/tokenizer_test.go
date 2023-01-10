package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBreakWord(t *testing.T) {
	t.Run("should tokenize", func(t *testing.T) {
		const word = "a1 b22\tc333\n\n\re4444  \f g"

		s := New()
		require.Equal(t, []string{
			"a1", "b22", "c333", "e4444", "g",
		},
			toStrings(s.BreakWord([]rune(word))),
		)
	})

	t.Run("should skip leading and trailing space", func(t *testing.T) {
		const word = " \na b\tc\n\n\re  \f g  "

		s := New()
		require.Equal(t, []string{
			"a", "b", "c", "e", "g",
		},
			toStrings(s.BreakWord([]rune(word))),
		)
	})
}

func toStrings(in [][]rune) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		out = append(out, string(s))
	}

	return out
}
