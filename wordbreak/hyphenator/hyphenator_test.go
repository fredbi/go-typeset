package hyphenator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestHyphenator(t *testing.T) {
	t.Run("should use en-US rules by default", func(t *testing.T) {
		h := New()

		require.Equal(t, []string{"an", "cil", "lary"}, h.BreakWord("ancillary"))
		require.Equal(t, []string{"as", "so", "ciate"}, h.BreakWord("associate"))
		require.Equal(t, []string{"associ-ates"}, h.BreakWord("associ-ates"))
		require.Equal(t, []string{"reci", "procity"}, h.BreakWord("reciprocity"))
		require.Equal(t, []string{"maj-estuous"}, h.BreakWord("maj-estuous"))
		require.Equal(t, []string{"ma", "jes", "tu", "ous"}, h.BreakWord("majestuous"))
	})

	t.Run("should use fr rules on option", func(t *testing.T) {
		h := New(WithLanguage("fr"))

		require.Equal(t, "patterns: hyph-fr.tex", h.String())
		require.Equal(t, []string{"connec", "ti", "vi", "té"}, h.BreakWord("connectivité"))
	})

	t.Run("should use de rules on option", func(t *testing.T) {
		h := New(WithLanguageTag(language.German))

		require.Equal(t, "German Hyphenation Patterns (Reformed Orthography, 2006) `dehyphn-x' 2019-04-04 (WL)}", h.String())
		require.Equal(t, []string{"Aus", "nah", "me"}, h.BreakWord("Ausnahme"))
	})

	t.Run("should use es rules on option", func(t *testing.T) {
		h := New(WithLanguageTag(language.Spanish))

		require.Equal(t, []string{"trans", "la", "to", "res"}, h.BreakWord("translatores"))
	})

	t.Run("should use en-GB rules on option", func(t *testing.T) {
		h := New(WithLanguageTag(language.BritishEnglish))

		require.Equal(t, []string{"uni", "ver", "sit", "ies"}, h.BreakWord("universities"))
	})

	t.Run("should fallback to some supported language", func(t *testing.T) {
		h := New(WithLanguageTag(language.Ukrainian))
		require.Equal(t, "patterns: ushyphmax.tex", h.String())
	})
}
