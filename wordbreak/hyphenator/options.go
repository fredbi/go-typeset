package hyphenator

import (
	"golang.org/x/text/language"
)

type (
	// Option to configure the hyphenator.
	Option func(*options)

	options struct {
		lang language.Tag
		// minRunesBeforeHyphen int
	}
)

// WithanguageTag specifies the language of the hyphenator, using
// a language tag from the standard library.
//
// The default language is "en-US".
func WithLanguageTag(tag language.Tag) Option {
	return func(o *options) {
		o.lang = tag
	}
}

// Withanguage specifies the language of the hyphenator, using
// a language string like "en-US", "en-GB", "es".
//
// The default language is "en-US".
func WithLanguage(lang string) Option {
	tag, _ := language.MatchStrings(langMatcher, lang)

	return func(o *options) {
		o.lang = tag
	}
}

func defaultOptions(opts []Option) *options {
	o := &options{
		lang: language.AmericanEnglish,
		// minRunesBeforeHyphen: 3,
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}
