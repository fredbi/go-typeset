//go:build profiler

package ansi

import (
	"testing"

	"github.com/pkg/profile"
)

func TestMemProfileStripANSI(t *testing.T) {
	const prof = "memprof"
	input := []rune(startInput + wordInput + endInput)

	defer profile.Start(
		profile.MemProfile,
		profile.ProfilePath(prof),
		profile.NoShutdownHook,
	).Stop()

	for n := 0; n < 10000; n++ {
		_, _, _ = StripANSIFromRunes(input)
	}
}

func TestCPUProfileStripANSI(t *testing.T) {
	const prof = "cpuprof"
	input := []rune(startInput + wordInput + endInput)

	defer profile.Start(
		profile.CPUProfile,
		profile.ProfilePath(prof),
		profile.NoShutdownHook,
	).Stop()

	for n := 0; n < 10000; n++ {
		_, _, _ = StripANSIFromRunes(input)
	}
}
