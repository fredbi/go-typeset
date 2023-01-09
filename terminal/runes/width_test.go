package runes

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

var runewidthtests = []struct {
	in     rune
	out    int
	eaout  int
	nseout int
}{
	{'世', 2, 2, 2},
	{'界', 2, 2, 2},
	{'ｾ', 1, 1, 1},
	{'ｶ', 1, 1, 1},
	{'ｲ', 1, 1, 1},
	{'☆', 1, 2, 2}, // double width in ambiguous
	{'☺', 1, 1, 2},
	{'☻', 1, 1, 2},
	{'♥', 1, 2, 2},
	{'♦', 1, 1, 2},
	{'♣', 1, 2, 2},
	{'♠', 1, 2, 2},
	{'♂', 1, 2, 2},
	{'♀', 1, 2, 2},
	{'♪', 1, 2, 2},
	{'♫', 1, 1, 2},
	{'☼', 1, 1, 2},
	{'↕', 1, 2, 2},
	{'‼', 1, 1, 2},
	{'↔', 1, 2, 2},
	{'\x00', 0, 0, 0},
	{'\x01', 0, 0, 0},
	{'\u0300', 0, 0, 0},
	{'\u2028', 0, 0, 0},
	{'\u2029', 0, 0, 0},
	{'a', 1, 1, 1}, // ASCII classified as "na" (narrow)
	{'⟦', 1, 1, 1}, // non-ASCII classified as "na" (narrow)
	{'👁', 1, 1, 2},
	{'ö', 1, 1, 1},
	{'ä', 1, 1, 1},
	{'æ', 1, 2, 2},
	{'ø', 1, 2, 2}, // ambiguous character defaulting to wide in East-Asian character sets
	{'Å', 1, 2, 2}, // ambiguous character defaulting to wide in East-Asian character sets
	{'å', 1, 1, 1}, // neutral character re East-Asian, defaults to narrow
}

func TestWidth(t *testing.T) {
	t.Run("With EastAsianWidth=false", func(t *testing.T) {
		t.Parallel()

		t.Run("should resolve rune widths", func(t *testing.T) {
			for _, toPin := range runewidthtests {
				testCase := toPin

				w := Width(testCase.in)
				require.Equal(t, runeWidthNonEastAsian(testCase.in), w)
				require.Equalf(t, testCase.out, w,
					"%[1]U: RuneWidth(%[1]q) = %d, want %d (EastAsianWidth=false)",
					testCase.in, w, testCase.out,
				)
			}
		})
	})

	t.Run("With EastAsianWidth=true", func(t *testing.T) {
		t.Parallel()

		t.Run("should resolve rune widths", func(t *testing.T) {
			for _, toPin := range runewidthtests {
				testCase := toPin

				w := Width(testCase.in, WithEastAsianWidth(true))
				require.Equal(t, runeWidthEastAsian(testCase.in, false), w)
				require.Equalf(t, testCase.eaout, w,
					"%[1]U: RuneWidth(%[1]q) = %d, want %d (EastAsianWidth=true, SkipStrictEmojiNeutral=false)",
					testCase.in, w, testCase.eaout,
				)
			}
		})
	})

	t.Run("With EastAsianWidth=true, SkipStrictEmojiNeutral=true", func(t *testing.T) {
		t.Parallel()

		t.Run("should resolve rune widths", func(t *testing.T) {
			for _, toPin := range runewidthtests {
				testCase := toPin
				w := Width(testCase.in, WithEastAsianWidth(true), WithSkipStrictEmojiNeutral(true))
				require.Equal(t, runeWidthEastAsian(testCase.in, true), w)
				require.Equalf(t, testCase.nseout, w,
					"%[1]U: RuneWidth(%[1]q) = %d, want %d (EastAsianWidth=true, SkipStrictEmojiNeutral=true)",
					testCase.in, w, testCase.nseout,
				)
			}
		})
	})

	t.Run("invalid rune should return 0 width", func(t *testing.T) {
		require.Equal(t, 0, Width(utf8.RuneError))

		const invalid = int32(0x0FFFFFFF)
		require.Equal(t, 0, Width(invalid))
	})

	t.Run("graphemes with multiple code points are not supported", func(t *testing.T) {
		const grapheme = "🏳️\u200d🌈"

		require.Equal(t, 14, len(grapheme))           // byte count
		require.Equal(t, 4, len([]rune(grapheme)))    // rune count
		require.NotEqual(t, 2, StringWidth(grapheme)) // should be 2 but we don't support this for now
	})
}

func BenchmarkWidth(b *testing.B) {
	// warm lookup cache
	_ = Width('A')

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		for _, testCase := range runewidthtests {
			_ = Width(testCase.in)
		}
	}
}

func BenchmarkWidths(b *testing.B) {
	// warm lookup cache
	_ = Width('A')
	str := []rune("string")

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		_ = Widths(str)
	}
}

func BenchmarkStringWidth(b *testing.B) {
	// warm lookup cache
	_ = Width('A')
	const str = "string"

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		_ = StringWidth(str)
	}
}

func TestAmbiguous(t *testing.T) {
	t.SkipNow()

	for _, runeRange := range ambiguous {
		for i := runeRange.first; i <= runeRange.last; i++ {
			t.Logf("%[1]U: %[1]c", i)
		}
	}
}
