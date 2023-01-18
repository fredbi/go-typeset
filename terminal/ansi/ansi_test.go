package ansi

import (
	"fmt"
	"testing"

	"github.com/fredbi/go-typeset/terminal/runes/runesio"
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

		rstripped := StripANSIFromRunes([]rune(input))
		require.Equal(t, input, string(rstripped.Text))
		require.Empty(t, rstripped.StartSequence)
		require.Empty(t, rstripped.StopSequence)
		require.Empty(t, string(rstripped.Remainder))
	})

	t.Run("strip ANSI should isolate string between start and end escape sequences", func(t *testing.T) {
		const input = startInput + wordInput + endInput

		rstripped := StripANSIFromRunes([]rune(input))
		require.Equal(t, wordInput, string(rstripped.Text))
		require.Equal(t, startInput, string(rstripped.StartSequence))
		require.Equal(t, endInput, string(rstripped.StopSequence))
		require.Empty(t, rstripped.Remainder)
	})

	t.Run("strip ANSI should isolate string when missing end escape sequence", func(t *testing.T) {
		const input = startInput + wordInput

		rstripped := StripANSIFromRunes([]rune(input))
		require.Equal(t, wordInput, string(rstripped.Text))
		require.Equal(t, startInput, string(rstripped.StartSequence))
		require.Empty(t, rstripped.StopSequence)
		require.Empty(t, rstripped.Remainder)
	})

	t.Run("strip ANSI should isolate string when missing start escape sequence", func(t *testing.T) {
		const input = wordInput + endInput

		rstripped := StripANSIFromRunes([]rune(input))
		require.Equal(t, wordInput, string(rstripped.Text))
		require.Empty(t, rstripped.StartSequence)
		require.Equal(t, endInput, string(rstripped.StopSequence))
		require.Empty(t, rstripped.Remainder)
	})

	t.Run("strip ANSI should isolate empty string when pure start/end escape sequence", func(t *testing.T) {
		const input = startInput + endInput

		rstripped := StripANSIFromRunes([]rune(input))
		require.Empty(t, rstripped.Text)
		require.Equal(t, startInput, string(rstripped.StartSequence))
		require.Equal(t, endInput, string(rstripped.StopSequence))
		require.Empty(t, rstripped.Remainder)
	})

	t.Run("strip ANSI should break down sequences when multiple start/end escape sequences", func(t *testing.T) {
		const input = startInput + startInput + wordInput + endInput + endInput

		rstripped := StripANSIFromRunes([]rune(input))
		require.Empty(t, rstripped.Text)
		require.Equal(t, startInput, string(rstripped.StartSequence))
		require.Empty(t, rstripped.StopSequence)
		require.NotEmpty(t, rstripped.Remainder)

		rstripped = StripANSIFromRunes(rstripped.Remainder)
		require.Equal(t, wordInput, string(rstripped.Text))
		require.Equal(t, startInput, string(rstripped.StartSequence))
		require.Equal(t, endInput, string(rstripped.StopSequence))
		require.NotEmpty(t, rstripped.Remainder)

		rstripped = StripANSIFromRunes(rstripped.Remainder)
		require.Empty(t, rstripped.Text)
		require.Empty(t, rstripped.StartSequence)
		require.Equal(t, endInput, string(rstripped.StopSequence))
		require.Empty(t, rstripped.Remainder)
	})

	t.Run("strip ANSI should find remainder string when multiple start/end escape sequences", func(t *testing.T) {
		const input = startInput + startInput + wordInput + endInput + endInput + startInput + wordInput + "2" + endInput

		rstripped := StripANSIFromRunes([]rune(input))
		require.Empty(t, rstripped.Text)
		require.Equal(t, startInput, string(rstripped.StartSequence))
		require.Empty(t, rstripped.StopSequence)
		require.NotEmpty(t, rstripped.Remainder)

		rstripped = StripANSIFromRunes(rstripped.Remainder)
		require.Equal(t, wordInput, string(rstripped.Text))
		require.Equal(t, startInput, string(rstripped.StartSequence))
		require.Equal(t, endInput, string(rstripped.StopSequence))
		require.NotEmpty(t, rstripped.Remainder)

		rstripped = StripANSIFromRunes(rstripped.Remainder)
		require.Empty(t, rstripped.Text)
		require.Empty(t, rstripped.StartSequence)
		require.Equal(t, endInput, string(rstripped.StopSequence))
		require.NotEmpty(t, rstripped.Remainder)

		rstripped = StripANSIFromRunes(rstripped.Remainder)
		require.Equal(t, wordInput+"2", string(rstripped.Text))
		require.Equal(t, startInput, string(rstripped.StartSequence))
		require.Equal(t, endInput, string(rstripped.StopSequence))
		require.Empty(t, rstripped.Remainder)
	})

	t.Run("strip ANSI recognize sequences with default numerical argument", func(t *testing.T) {
		stripped := StripANSIFromRunes([]rune("\033[1mABC\033[1;2~\033[K\033[m"))
		require.Equal(t, "ABC\033[1;2~\033[K", string(stripped.Text)) // this sequence, not recognized as either Start or Stop, is left in the text
		require.Equal(t, "\x1b[1m", string(stripped.StartSequence))
		require.Equal(t, "\x1b[m", string(stripped.StopSequence))
		require.Empty(t, stripped.Remainder)
	})

	t.Run("strip ANSI recognize sequences with ':' as an arguments separator", func(t *testing.T) {
		stripped := StripANSIFromRunes([]rune("\033[1:2mABC\033[m\033[1:2;;~"))
		require.Equal(t, "ABC", string(stripped.Text))
		require.Equal(t, "\x1b[1:2m", string(stripped.StartSequence))
		require.Equal(t, "\x1b[m", string(stripped.StopSequence))
		require.NotEmpty(t, stripped.Remainder)

		stripped = StripANSIFromRunes(stripped.Remainder)
		require.Equal(t, "\x1b[1:2;;~", string(stripped.Text))
		require.Empty(t, stripped.StartSequence)
		require.Empty(t, stripped.StopSequence)
		require.Empty(t, stripped.Remainder)
	})

	t.Run("strip ANSI recognize sequences with more than 2 numerical arguments", func(t *testing.T) {
		stripped := StripANSIFromRunes([]rune("\033[1;2;3;4mABC\033[1;2~\033[1:2:3:4m"))
		require.Equal(t, "ABC\033[1;2~", string(stripped.Text))
		require.Equal(t, "\x1b[1;2;3;4m", string(stripped.StartSequence))
		require.Empty(t, stripped.StopSequence)
		require.NotEmpty(t, stripped.Remainder)

		stripped = StripANSIFromRunes(stripped.Remainder)
		require.Empty(t, stripped.Text)
		require.Equal(t, "\x1b[1:2:3:4m", string(stripped.StartSequence))
		require.Empty(t, stripped.StopSequence)
		require.Empty(t, stripped.Remainder)
	})

	t.Run("strip ANSI should break down whe first unrecognized sequence after start sequence", func(t *testing.T) {
		stripped := StripANSIFromRunes([]rune("\033[4m\033[1;2~\033[m"))
		require.Equal(t, "\033[1;2~", string(stripped.Text))
		require.Equal(t, "\033[4m", string(stripped.StartSequence))
		require.Empty(t, stripped.StopSequence)
		require.NotEmpty(t, stripped.Remainder)

		stripped = StripANSIFromRunes(stripped.Remainder)
		require.Empty(t, stripped.Text)
		require.Empty(t, stripped.StartSequence)
		require.Equal(t, "\033[m", string(stripped.StopSequence))
		require.Empty(t, stripped.Remainder)
	})

	t.Run("strip ANSI should break down when first unrecognized sequence after start sequence", func(t *testing.T) {
		stripped := StripANSIFromRunes([]rune("\033[1;2~\033[4m"))
		require.Equal(t, "\033[1;2~", string(stripped.Text))
		require.Empty(t, stripped.StartSequence)
		require.Empty(t, stripped.StopSequence)
		require.NotEmpty(t, stripped.Remainder)

		stripped = StripANSIFromRunes(stripped.Remainder)
		require.Empty(t, stripped.Text)
		require.Equal(t, "\033[4m", string(stripped.StartSequence))
		require.Empty(t, stripped.StopSequence)
		require.Empty(t, stripped.Remainder)
	})

	t.Run("strip ANSI should break down as text when ANSI escape sequence is incomplete", func(t *testing.T) {
		stripped := StripANSIFromRunes([]rune("\033[1;2m\033X\033[m")) // we have ESC, but not [
		require.Equal(t, "\033X", string(stripped.Text))
		require.Equal(t, "\033[1;2m", string(stripped.StartSequence))
		require.Equal(t, "\033[m", string(stripped.StopSequence))
		require.Empty(t, stripped.Remainder)
	})

	t.Run("strip ANSI should break down as text when ANSI escape sequence is truncated", func(t *testing.T) {
		stripped := StripANSIFromRunes([]rune("\033[1;2m\033\033"))
		require.Equal(t, "\033\033", string(stripped.Text))
		require.Equal(t, "\033[1;2m", string(stripped.StartSequence))
		require.Empty(t, stripped.StopSequence)
		require.Empty(t, stripped.Remainder)
	})

	t.Run("strip ANSI should break down as text when ANSI escape sequence is truncated (2)", func(t *testing.T) {
		stripped := StripANSIFromRunes([]rune("\033"))
		require.Equal(t, "\033", string(stripped.Text))
		require.Empty(t, stripped.StartSequence)
		require.Empty(t, stripped.StopSequence)
		require.Empty(t, stripped.Remainder)
	})

	t.Run("strip empty runes", func(t *testing.T) {
		stripped := StripANSIFromRunes([]rune{})
		require.Empty(t, stripped.Text)
		require.Empty(t, stripped.StartSequence)
		require.Empty(t, stripped.StopSequence)
		require.Empty(t, stripped.Remainder)
	})
}

func TestDecodeANSISequence(t *testing.T) {
	// test a variety of legit ANSI sequences (test sequence scanner)

	for _, toPin := range seqTestCases() {
		testCase := toPin

		t.Run(fmt.Sprintf("should detect start sequence %s", testCase.Title), func(t *testing.T) {
			t.Parallel()

			rdr := runesio.NewSliceReader(testCase.Input)
			seq, seqKind := decodeANSISequence(rdr)

			if testCase.ExpectedStart != nil {
				require.Equal(t, testCase.ExpectedStart, seq)
				require.Equal(t, startSequence, seqKind)
			}

			if testCase.ExpectedStop != nil {
				require.Equal(t, testCase.ExpectedStop, seq)
				require.Equal(t, stopSequence, seqKind)
			}

			if testCase.ExpectedOther != nil {
				require.Equal(t, testCase.ExpectedOther, seq)
				require.Equal(t, otherSequence, seqKind)
			}

			if testCase.ExpectedRemainder != nil {
				remainder := rdr.Runes()
				require.Equal(t, testCase.ExpectedRemainder, remainder)
			}
		})
	}
}

type seqTestCase struct {
	Title             string
	Input             []rune
	ExpectedStart     []rune
	ExpectedStop      []rune
	ExpectedOther     []rune
	ExpectedRemainder []rune
}

func seqTestCases() []seqTestCase {
	return []seqTestCase{
		{
			Title:             "with SGR start",
			Input:             []rune("\033[12mTEXT"),
			ExpectedStart:     []rune("\033[12m"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with SGR start, extra args",
			Input:             []rune("\033[12;1;2;3mTEXT"),
			ExpectedStart:     []rune("\033[12;1;2;3m"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with [s start",
			Input:             []rune("\033[12sTEXT"),
			ExpectedStart:     []rune("\033[12s"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with [s start, no arg",
			Input:             []rune("\033[sTEXT"),
			ExpectedStart:     []rune("\033[s"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with [h start",
			Input:             []rune("\033[12hTEXT"),
			ExpectedStart:     []rune("\033[12h"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with [h start, no arg",
			Input:             []rune("\033[hTEXT"),
			ExpectedStart:     []rune("\033[h"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with SGR stop",
			Input:             []rune("\033[0mTEXT"),
			ExpectedStop:      []rune("\033[0m"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with SGR stop (multiple zeros)",
			Input:             []rune("\033[00mTEXT"),
			ExpectedStop:      []rune("\033[00m"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with SGR stop, no arg",
			Input:             []rune("\033[mTEXT"),
			ExpectedStop:      []rune("\033[m"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with [u stop",
			Input:             []rune("\033[12uTEXT"),
			ExpectedStop:      []rune("\033[12u"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with [u stop, no arg",
			Input:             []rune("\033[uTEXT"),
			ExpectedStop:      []rune("\033[u"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with [l stop",
			Input:             []rune("\033[12lTEXT"),
			ExpectedStop:      []rune("\033[12l"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with [l stop, no arg",
			Input:             []rune("\033[lTEXT"),
			ExpectedStop:      []rune("\033[l"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with other sequence",
			Input:             []rune("\033[1;2KTEXT"),
			ExpectedOther:     []rune("\033[1;2K"),
			ExpectedRemainder: []rune("TEXT"),
		},
		{
			Title:             "with non-ANSI sequence",
			Input:             []rune("\033#4TEXT"),
			ExpectedRemainder: []rune("\033#4TEXT"),
		},
		{
			Title:         "with incomplete sequence (no key code)",
			Input:         []rune("\033["),
			ExpectedOther: []rune("\033["),
		},
		{
			Title:         "with incomplete sequence (no arg)",
			Input:         []rune("\033[;"),
			ExpectedOther: []rune("\033[;"),
		},
		{
			Title:             "with incomplete sequence (no CSI)",
			Input:             []rune("\033"),
			ExpectedRemainder: []rune("\033"),
		},
		{
			Title:             "with incomplete sequence (no CSI) (2)",
			Input:             []rune("\033\033\033"),
			ExpectedRemainder: []rune("\033\033\033"),
		},
	}
}

func TestStripToken(t *testing.T) {
	const token = "Lorem ipsum dolor \x1b[31;1;4msit amet,"
	stripped := StripToken([]rune(token))
	require.Len(t, stripped, 2)
	require.Equal(t, StrippedToken{
		Text:      []rune("Lorem ipsum dolor "),
		Remainder: []rune("\033[31;1;4msit amet,"),
	}, stripped[0])
	require.Equal(t, StrippedToken{
		Text:          []rune("sit amet,"),
		StartSequence: []rune("\033[31;1;4m"),
		Remainder:     []rune{},
	}, stripped[1])
}
