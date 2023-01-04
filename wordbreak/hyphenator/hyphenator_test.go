package hyphenator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

const superLongWord = `Honorificabilitudinitatibus`

func TestHyphenator(t *testing.T) {
	t.Parallel()

	t.Run("should use en-US rules by default", func(t *testing.T) {
		t.Parallel()
		h := New()

		require.Equal(t, toRunes([]string{"an", "cil", "lary"}), h.BreakWordString("ancillary"))
		require.Equal(t, toRunes([]string{"as", "so", "ciate"}), h.BreakWordString("associate"))
		require.Equal(t, toRunes([]string{"associ-ates"}), h.BreakWordString("associ-ates"))
		require.Equal(t, toRunes([]string{"reci", "procity"}), h.BreakWordString("reciprocity"))
		require.Equal(t, toRunes([]string{"maj-estuous"}), h.BreakWordString("maj-estuous"))
		require.Equal(t, toRunes([]string{"ma", "jes", "tu", "ous"}), h.BreakWordString("majestuous"))
		require.Equal(t, toRunes([]string{"daugh", "ters"}), h.BreakWordString("daughters"))
		require.Equal(t, toRunes([]string{"sub", "di", "vi", "sion"}), h.BreakWordString("subdivision"))
		require.Equal(t, toRunes([]string{"dis", "ci", "plines"}), h.BreakWordString("disciplines"))
		require.Equal(t, toRunes([]string{"phil", "an", "thropic"}), h.BreakWordString("philanthropic"))
	})

	t.Run("should abide by min length rules", func(t *testing.T) {
		t.Parallel()
		h := New()

		require.Equal(t, toRunes([]string{"ex", "am", "ple"}), h.BreakWordString("example"))
		require.Equal(t, toRunes([]string{"king"}), h.BreakWordString("king"))
		require.Equal(t, toRunes([]string{"daugh", "ter"}), h.BreakWordString("daughter"))
		require.Equal(t, toRunes([]string{"daugh", "ters"}), h.BreakWordString("daughters"))
	})

	t.Run("pattern search should not be case sensitive", func(t *testing.T) {
		t.Parallel()
		h := New()

		require.Equal(t, toRunes([]string{"Ex", "am", "ple"}), h.BreakWordString("Example"))
	})

	t.Run("should use fr rules on option", func(t *testing.T) {
		t.Parallel()
		h := New(WithLanguage("fr"))

		require.Equal(t, "patterns: hyph-fr.tex", h.String())
		require.Equal(t, toRunes([]string{"connec", "ti", "vi", "té"}), h.BreakWordString("connectivité"))
	})

	t.Run("should use de rules on option", func(t *testing.T) {
		t.Parallel()
		h := New(WithLanguageTag(language.German))

		require.Equal(t, "German Hyphenation Patterns (Reformed Orthography, 2006) `dehyphn-x' 2019-04-04 (WL)}", h.String())
		require.Equal(t, toRunes([]string{"Aus", "nah", "me"}), h.BreakWordString("Ausnahme"))
		require.Equal(t, toRunes([]string{"Uni", "ver", "si", "täts", "stadt"}), h.BreakWordString("Universitätsstadt"))

	})

	t.Run("should use es rules on option", func(t *testing.T) {
		t.Parallel()
		h := New(WithLanguageTag(language.Spanish))

		require.Equal(t, toRunes([]string{"trans", "la", "to", "res"}), h.BreakWordString("translatores"))
	})

	t.Run("should use en-GB rules on option", func(t *testing.T) {
		t.Parallel()
		h := New(WithLanguageTag(language.BritishEnglish))

		require.Equal(t, toRunes([]string{"uni", "ver", "sit", "ies"}), h.BreakWordString("universities"))
	})

	t.Run("should fallback to some supported language", func(t *testing.T) {
		t.Parallel()
		h := New(WithLanguageTag(language.Ukrainian))

		require.Equal(t, "patterns: ushyphmax.tex", h.String())
	})

	t.Run("should break super long word", func(t *testing.T) {
		t.Parallel()
		h := New()

		require.Equal(t, toRunes([]string{"Hon", "ori", "fi", "ca", "bil", "itu", "dini", "tat", "ibus"}), h.BreakWordString(superLongWord))
	})

	t.Run("should not hyphen on shorter alternatives", func(t *testing.T) {
		t.Parallel()
		h := New(WithMinLength(6), WithMinLeft(6), WithMinRight(2))

		res := h.BreakWordString(superLongWord)
		require.Equalf(t, toRunes([]string{"Honori", "ficabil", "itudini", "tatibus"}), res,
			"unexpected: %v", toStrings(res),
		)
	})

	t.Run("should not hyphen shorter words", func(t *testing.T) {
		t.Parallel()
		h := New(WithMinLength(7), WithMinLeft(6), WithMinRight(2))

		res := h.BreakWordString("example")
		require.Equalf(t, toRunes([]string{"example"}), res,
			"unexpected: %v", toStrings(res),
		)
	})

	t.Run("should hyphen unknown words (and successfully extend positions)", func(t *testing.T) {
		t.Parallel()
		h := New()

		res := h.BreakWordString(strings.Repeat(superLongWord, 2))
		require.Equalf(t, toRunes([]string{
			"Hon", "ori", "fi", "ca", "bil", "itu", "dini", "tat", "ibusHon",
			"ori", "fi", "ca", "bil", "itu", "dini", "tat", "ibus",
		}), res,
			"unexpected: %v", toStrings(res),
		)
	})
}

func TestSplitWord(t *testing.T) {
	t.Parallel()

	t.Run("should split according to hyphens", func(t *testing.T) {
		require.Equal(t, toRunes([]string{"XY"}), SplitWord([]rune("XY")))
		require.Equal(t, toRunes([]string{"X", "-", "Y"}), SplitWord([]rune("X-Y")))
		require.Equal(t, toRunes([]string{"X", "-", "Y", "-"}), SplitWord([]rune("X-Y-")))
		require.Equal(t, toRunes([]string{"-", "X", "-", "Y"}), SplitWord([]rune("-X-Y")))
		require.Equal(t, toRunes([]string{"X", "-", "-", "Y"}), SplitWord([]rune("X--Y")))
		require.Equal(t, toRunes([]string{"X", "-", "-"}), SplitWord([]rune("X--")))
		require.Equal(t, toRunes([]string{"-", "-", "X"}), SplitWord([]rune("--X")))
	})
}

func TestIsHyphen(t *testing.T) {
	t.Parallel()

	require.True(t, IsHyphen([]rune("-")))
	require.False(t, IsHyphen([]rune("_")))
	require.False(t, IsHyphen([]rune("foo")))
	require.False(t, IsHyphen([]rune("")))
}

func toRunes(in []string) [][]rune {
	out := make([][]rune, 0, len(in))
	for _, s := range in {
		out = append(out, []rune(s))
	}

	return out
}

func toStrings(in [][]rune) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		out = append(out, string(s))
	}

	return out
}
