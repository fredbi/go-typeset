package ansi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	startInput = "\033[43;30m"
	endInput   = "\033[00m"
	wordInput  = "Česká řeřicha"
)

func TestStripANSI(t *testing.T) {
	t.Parallel()

	t.Run("strip ANSI should leave non-escaped string unchanged", func(t *testing.T) {
		const input = "ABC"
		stripped, start, end, remain := StripANSI(input)

		require.Equal(t, input, stripped)
		require.Empty(t, start)
		require.Empty(t, end)
		require.Empty(t, remain)

		rstripped, rstart, rend, remainder := StripANSIFromRunes([]rune(input))
		require.Equal(t, stripped, string(rstripped))
		require.Empty(t, rstart)
		require.Empty(t, rend)
		require.Empty(t, string(remainder))
	})

	t.Run("strip ANSI should isolate string between start and end escape sequences", func(t *testing.T) {
		const input = startInput + wordInput + endInput

		stripped, start, end, remain := StripANSI(input)
		require.Equal(t, wordInput, stripped)
		require.Equal(t, startInput, start)
		require.Equal(t, endInput, end)
		require.Empty(t, remain)

		rstripped, rstart, rend, remainder := StripANSIFromRunes([]rune(input))
		require.Equal(t, stripped, string(rstripped))
		require.Equal(t, start, string(rstart))
		require.Equal(t, end, string(rend))
		require.Empty(t, remainder)
	})

	t.Run("strip ANSI should isolate string when missing end escape sequence", func(t *testing.T) {
		const input = startInput + wordInput

		stripped, start, end, remain := StripANSI(input)
		require.Equal(t, wordInput, stripped)
		require.Equal(t, startInput, start)
		require.Empty(t, end)
		require.Empty(t, remain)

		rstripped, rstart, rend, remainder := StripANSIFromRunes([]rune(input))
		require.Equal(t, stripped, string(rstripped))
		require.Equal(t, start, string(rstart))
		require.Equal(t, end, string(rend))
		require.Empty(t, string(remainder))
	})

	t.Run("strip ANSI should isolate string when missing start escape sequence", func(t *testing.T) {
		const input = wordInput + endInput

		stripped, start, end, remain := StripANSI(input)
		require.Equal(t, wordInput, stripped)
		require.Empty(t, start)
		require.Equal(t, endInput, end)
		require.Empty(t, remain)

		rstripped, rstart, rend, remainder := StripANSIFromRunes([]rune(input))
		require.Equal(t, stripped, string(rstripped))
		require.Equal(t, start, string(rstart))
		require.Equal(t, end, string(rend))
		require.Empty(t, remainder)
	})

	t.Run("strip ANSI should isolate empty string when pure start/end escape sequence", func(t *testing.T) {
		const input = startInput + endInput

		stripped, start, end, remain := StripANSI(input)
		require.Empty(t, stripped)
		require.Equal(t, startInput+endInput, start)
		require.Empty(t, end)
		require.Empty(t, remain)

		rstripped, rstart, rend, remainder := StripANSIFromRunes([]rune(input))
		require.Equal(t, stripped, string(rstripped))
		require.Equal(t, start, string(rstart))
		require.Equal(t, end, string(rend))
		require.Empty(t, remainder)
	})

	t.Run("strip ANSI should isolate string when multiple start/end escape sequences", func(t *testing.T) {
		const input = startInput + startInput + wordInput + endInput + endInput

		stripped, start, end, remain := StripANSI(input)
		require.Equal(t, wordInput, stripped)
		require.Equal(t, startInput+startInput, start)
		require.Equal(t, endInput+endInput, end)
		require.Empty(t, remain)

		rstripped, rstart, rend, remainder := StripANSIFromRunes([]rune(input))
		require.Equal(t, stripped, string(rstripped))
		require.Equal(t, start, string(rstart))
		require.Equal(t, end, string(rend))
		require.Empty(t, remainder)
	})

	t.Run("strip ANSI should find remainder string when multiple start/end escape sequences", func(t *testing.T) {
		const input = startInput + startInput + wordInput + endInput + endInput + startInput + wordInput + "2" + endInput

		stripped, start, end, remain := StripANSI(input)
		require.Equal(t, wordInput, stripped)
		require.Equal(t, startInput+startInput, start)
		require.Equal(t, endInput+endInput+startInput, end)
		require.Equal(t, wordInput+"2"+endInput, remain)

		rstripped, rstart, rend, remainder := StripANSIFromRunes([]rune(input))
		require.Equal(t, wordInput, string(rstripped))
		require.Equal(t, startInput+startInput, string(rstart))
		require.Equal(t, endInput+endInput+startInput, string(rend))
		require.Equal(t, wordInput+"2"+endInput, string(remainder))

		stripped, start, end, remain = StripANSI(remain)
		require.Equal(t, wordInput+"2", stripped)
		require.Empty(t, start)
		require.Equal(t, endInput, end)
		require.Empty(t, remain)
	})
}

// test a variety of legit ANSI sequences (test scanner)
func TestDecodeANSISequence(t *testing.T) {
	t.Run("strip ANSI recognize sequences with default numerical argument", func(t *testing.T) {
		stripped, start, end, remainder := StripANSIFromRunes([]rune("\033[mABC\033[1;2~"))
		require.Equal(t, "ABC", string(stripped))
		require.Equal(t, "\x1b[m", string(start))
		require.Equal(t, "\x1b[1;2~", string(end))
		require.Empty(t, remainder)
	})

	t.Run("strip ANSI recognize sequences with ':' as an arguments separator", func(t *testing.T) {
		stripped, start, end, remainder := StripANSIFromRunes([]rune("\033[1:2mABC\033[1:2;;~"))
		require.Equal(t, "ABC", string(stripped))
		require.Equal(t, "\x1b[1:2m", string(start))
		require.Equal(t, "\x1b[1:2;;~", string(end))
		require.Empty(t, remainder)
	})

	t.Run("strip ANSI recognize sequences with more than 2 numerical arguments", func(t *testing.T) {
		stripped, start, end, remainder := StripANSIFromRunes([]rune("\033[1;2;3;4mABC\033[1:2:3:4~"))
		require.Equal(t, "ABC", string(stripped))
		require.Equal(t, "\x1b[1;2;3;4m", string(start))
		require.Equal(t, "\x1b[1:2:3:4~", string(end))
		require.Empty(t, remainder)
	})

}
