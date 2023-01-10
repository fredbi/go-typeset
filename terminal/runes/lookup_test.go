package runes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsAmbiguous(t *testing.T) {
	require.True(t, IsAmbiguous('ø'))
	require.True(t, IsAmbiguous('Å'))
	require.True(t, IsAmbiguous('æ'))
	require.False(t, IsAmbiguous('å'))
}
