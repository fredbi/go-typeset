package hyphenator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSupportedPatterns(t *testing.T) {
	supported, err := SupportedPatterns()
	require.NoError(t, err)

	require.Equal(t, []string{
		"hyph-de-1996.tex",
		"hyph-en-gb.tex",
		"hyph-es.tex",
		"hyph-fr.tex",
		"ushyphmax.tex",
	}, supported,
	)
}

func TestLoadPatterns(t *testing.T) {
	supported, err := SupportedPatterns()
	require.NoError(t, err)

	for _, patterns := range supported {
		dict, err := LoadPatterns(patterns)
		require.NoError(t, err)
		require.NotNil(t, dict)
	}
}
