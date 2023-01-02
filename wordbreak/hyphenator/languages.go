package hyphenator

import (
	"golang.org/x/text/language"
)

var (
	supportedLanguages = []language.Tag{
		language.AmericanEnglish, // en-US
		language.BritishEnglish,  // en-GB
		language.English,         // en
		language.French,          // fr
		language.CanadianFrench,
		language.German,  // de
		language.Spanish, // es
		language.EuropeanSpanish,
		language.LatinAmericanSpanish,
	}

	langMatcher = language.NewMatcher(supportedLanguages)
)

func langToPattern(tag language.Tag) string {
	matched, _, _ := langMatcher.Match(tag)

	switch matched {
	case language.BritishEnglish, language.English:
		return "hyph-en-gb.tex"
	case language.Spanish, language.EuropeanSpanish, language.LatinAmericanSpanish:
		return "hyph-es.tex"
	case language.French, language.CanadianFrench:
		return "hyph-fr.tex"
	case language.German:
		return "hyph-de-1996.tex"
	case language.AmericanEnglish:
		fallthrough
	default:
		return "ushyphmax.tex"
	}
}
