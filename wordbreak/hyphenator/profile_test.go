//go:build profiler

package hyphenator

import (
	"testing"

	"github.com/pkg/profile"
	"github.com/stretchr/testify/require"
)

func TestMemProfileHyphenator(t *testing.T) {
	const prof = "memprof"
	superLongRunes := []rune(superLongWord)

	h := New()
	// pin the cache
	_ = h.BreakWord(superLongRunes)

	defer profile.Start(
		profile.MemProfile,
		profile.ProfilePath(prof),
		profile.NoShutdownHook,
	).Stop()

	for n := 0; n < 10000; n++ {
		_ = h.BreakWord(superLongRunes)
	}
}

func TestCPUProfileHyphenator(t *testing.T) {
	const prof = "cpuprof"
	superLongRunes := []rune(superLongWord)

	h := New()
	// pin the cache
	_ = h.BreakWord(superLongRunes)

	defer profile.Start(
		profile.CPUProfile,
		profile.ProfilePath(prof),
		profile.NoShutdownHook,
	).Stop()

	for n := 0; n < 10000; n++ {
		_ = h.BreakWord(superLongRunes)
	}
}

func TestMemProfileLoadPattern(t *testing.T) {
	const prof = "mempatprof"

	supported, err := SupportedPatterns()
	require.NoError(t, err)
	patterns := supported[0]

	defer profile.Start(
		profile.MemProfile,
		profile.ProfilePath(prof),
		profile.NoShutdownHook,
	).Stop()

	dict, err := LoadPatterns(patterns)
	require.NoError(t, err)
	require.NotNil(t, dict)
}

func TestMemProfileSplitWord(t *testing.T) {
	const prof = "memsplitprof"
	var superLongSplit = []rune(`Hono-rifica-bi-litu-dini-ta-ti-bus`)

	defer profile.Start(
		profile.MemProfile,
		profile.ProfilePath(prof),
		profile.NoShutdownHook,
	).Stop()

	for n := 0; n < 10000; n++ {
		_ = SplitWord(superLongSplit)
	}
}
