package runesio

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunesWriter(t *testing.T) {
	w := new(bytes.Buffer)

	writer := AsRunesWriter(w)

	t.Run("should write a single rune", func(t *testing.T) {
		res, err := writer.WriteRune('a')
		require.NoError(t, err)
		require.Equal(t, 1, res)
		require.Equal(t, "a", w.String())
	})

	t.Run("should write several runes", func(t *testing.T) {
		res, err := writer.WriteRunes([]rune{'a', 'b', 'c'})
		require.Equal(t, "aabc", w.String())
		require.NoError(t, err)
		require.Equal(t, 3, res)
	})

	t.Run("should report the size in bytes", func(t *testing.T) {
		res, err := writer.WriteRunes([]rune{'é', '🏈'})
		require.Equal(t, "aabcé🏈", w.String())
		require.NoError(t, err)
		require.Equal(t, 6, res)
	})
}
