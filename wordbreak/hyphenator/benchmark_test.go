package hyphenator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkHyphenatorLong(b *testing.B) {
	superLongRunes := []rune(superLongWord)

	h := New()
	// pin the cache
	_ = h.BreakWord([]rune{})

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		_ = h.BreakWord(superLongRunes)
	}
}

func BenchmarkHyphenatorShort(b *testing.B) {
	regularRunes := []rune("example")

	h := New()
	// pin the cache
	_ = h.BreakWord([]rune{})

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		_ = h.BreakWord(regularRunes)
	}
}

func BenchmarkHyphenatorString(b *testing.B) {
	h := New()
	// pin the cache
	_ = h.BreakWordString("pin")

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		_ = h.BreakWordString(superLongWord)
	}
}

func BenchmarkLoadPatterns(b *testing.B) {
	supported, err := SupportedPatterns()
	require.NoError(b, err)
	patterns := supported[0]

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		dict, err := LoadPatterns(patterns)
		require.NoError(b, err)
		require.NotNil(b, dict)
	}
}
