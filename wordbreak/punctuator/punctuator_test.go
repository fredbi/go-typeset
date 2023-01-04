package punctuator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPunctuator(t *testing.T) {
	t.Parallel()

	s := New()
	t.Run("should break along punctuation signs", func(t *testing.T) {
		const word = `a/b/c|d-e_foo\g-`
		require.Equal(t, toRunes([]string{
			"a", "/", "b", "/", "c", "|", "d", "-", "e", "_", "foo", `\`, "g", "-",
		}),
			s.BreakWord([]rune(word)),
		)
	})

	t.Run("should handle edge cases", func(t *testing.T) {
		t.Run("with heading punct", func(t *testing.T) {
			var word = []rune(`,foo`)
			require.Equal(t, toRunes([]string{
				",", "foo",
			}),
				s.BreakWord(word),
			)
		})
		t.Run("with trailing punct", func(t *testing.T) {
			var word = []rune(`foo,`)
			require.Equal(t, toRunes([]string{
				"foo", ",",
			}),
				s.BreakWord(word),
			)
		})
		t.Run("with duplicate punct", func(t *testing.T) {
			var word = []rune(`foo,,`)
			require.Equal(t, toRunes([]string{
				"foo", ",", ",",
			}),
				s.BreakWord(word),
			)
		})

		t.Run("with string input", func(t *testing.T) {
			const word = `foo,,`
			require.Equal(t, toRunes([]string{
				"foo", ",", ",",
			}),
				s.BreakWordString(word),
			)
		})
	})
}

func TestIsBreakRune(t *testing.T) {
	t.Parallel()

	for _, r := range []rune{
		'/',
		'_',
		'-',
		'&',
		'!',
		'|',
	} {
		require.True(t, punctSplitFunc(r))
		require.True(t, IsPunctuation([]rune{r}))
	}
}

func TestIsNotPunctuation(t *testing.T) {
	t.Parallel()

	require.False(t, IsPunctuation([]rune("aa")))
	require.False(t, IsPunctuation([]rune(",,")))
	require.False(t, IsPunctuation([]rune("")))
}

func toRunes(in []string) [][]rune {
	out := make([][]rune, 0, len(in))
	for _, s := range in {
		out = append(out, []rune(s))
	}

	return out
}
