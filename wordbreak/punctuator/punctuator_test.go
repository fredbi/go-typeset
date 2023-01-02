package punctuator

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

func TestPunctuator(t *testing.T) {
	t.Run("should break along punctuation signs", func(t *testing.T) {
		const word = `a/b/c|d-e_f\g-`

		s := New()
		require.Equal(t, []string{
			"a/", "b/", "c|", "d-", "e_", `f\`, "g-",
		},
			s.BreakWord(word),
		)
	})
}

func TestRunes(t *testing.T) {
	t.SkipNow()

	for _, r := range []rune{
		'/',
		'_',
		'-',
		'&',
		'!',
		'|',
	} {
		t.Logf("isPunct(%v)=%t [%t]", r, unicode.IsPunct(r), punctSplitFunc(r))
	}
}
