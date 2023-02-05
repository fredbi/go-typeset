package hyphenator

import (
	"golang.org/x/text/language"
)

type (
	// Option to configure the hyphenator.
	Option func(*options)

	options struct {
		lang      language.Tag
		minLength int
		minLeft   int
		minRight  int
	}
)

// WithanguageTag specifies the language of the hyphenator, using
// a language tag from the standard library.
//
// The default language is "language.AmericanEnglish".
//
// Unsupported languages are matched to sensible defaults, using a language.Matcher.
func WithLanguageTag(tag language.Tag) Option {
	return func(o *options) {
		o.lang = tag
	}
}

// Withanguage specifies the language of the hyphenator, using
// a language string like "en-US", "en-GB", "es".
//
// The default language is "en-US".
//
// Unsupported languages are matched to sensible defaults, using a language.Matcher.
func WithLanguage(lang string) Option {
	tag, _ := language.MatchStrings(langMatcher, lang)

	return func(o *options) {
		o.lang = tag
	}
}

// WithMinLength configures the minimum length (in runes) of a word to be eligible to hyphenation.
//
// The default is 4.
func WithMinLength(minLength int) Option {
	return func(o *options) {
		o.minLength = minLength
	}
}

// WithMinLeft configures the minimum length of a word part before an hyphenation point.
//
// The default is 2, meaning that the hyphenator will never break words leaving a single rune to the left.
func WithMinLeft(minLeft int) Option {
	return func(o *options) {
		o.minLeft = minLeft
	}
}

// WithMinRight configures the minimum length of a word part after an hyphenation point.
//
// The default is 2, meaning that the hyphenator will never break words leaving a single rune to the right.
func WithMinRight(minRight int) Option {
	return func(o *options) {
		o.minRight = minRight
	}
}

func defaultOptions(opts []Option) *options {
	o := &options{
		lang:      language.AmericanEnglish,
		minLength: 4,
		minLeft:   2,
		minRight:  2,
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}
