package runes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldsFunc(t *testing.T) {
	splitter := func(r rune) bool { return r == '|' }

	t.Run("should split text as runes", func(t *testing.T) {
		const word = "a1|b22|c333||||e4444|||g"

		require.Equal(t, []string{
			"a1", "b22", "c333", "e4444", "g",
		},
			toStrings(FieldsFunc([]rune(word), splitter)),
		)
	})

	t.Run("should skip leading and trailing separators", func(t *testing.T) {
		const word = "||a|b|c|||e||g|||"

		require.Equal(t, []string{
			"a", "b", "c", "e", "g",
		},
			toStrings(FieldsFunc([]rune(word), splitter)),
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
