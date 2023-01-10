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

/* TODO: test other sequences (test scanner)
func TestDecodeANSISequence(t *testing.T) {
	// var input = []rune(startInput + startInput + wordInput + endInput + endInput + startInput + endInput)
	var input = []rune(startInput + wordInput + endInput)
	stripped, start, end, _ := StripANSIFromRunes(input)
	t.Logf("DEBUG: %q, %q, %q", string(stripped), string(start), string(end))
	//func StripANSIFromRunes2(rns []rune) ([]rune, []rune, []rune) {

}
*/
