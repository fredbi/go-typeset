package runesio

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunesWriter(t *testing.T) {
	w := new(bytes.Buffer)

	writer := AsRunesWriter(w)

	res, err := writer.WriteRune('a')
	require.NoError(t, err)
	require.Equal(t, 1, res)
	require.Equal(t, "a", w.String())
	res, err = writer.WriteRunes([]rune{'a', 'b', 'c'})
	require.Equal(t, "aabc", w.String())
	require.NoError(t, err)
	require.Equal(t, 1, res)
	res, err = writer.WriteRunes([]rune{'é', '🏈'})
	require.Equal(t, "aabcé🏈", w.String())
	require.NoError(t, err)
	require.Equal(t, 1, res)
}
