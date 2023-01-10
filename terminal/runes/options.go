package runes

type (
	// Option to set expectations on rune width computation behavior.
	Option func(*options)

	options struct {
		EastAsian                  bool
		SkipStrictEmojiNeutral     bool
		DefaultAsianAmbiguousWidth int
	}
)

var (
	defaultOptions = &options{
		EastAsian:                  false,
		SkipStrictEmojiNeutral:     false,
		DefaultAsianAmbiguousWidth: 2,
	}
)

// WithEastAsianWidth should be set true when displaying Asian characters.
//
// The difference is that "ambiguous" runes that appear in different character sets
// and are normally with width 1 are not reported with width 2.
//
// Example: 'ø' is reported with width 1 by default (EastAsianWidth disabled)
// and with width 2 when EastAsianWidth is reported.
func WithEastAsianWidth(enabled bool) Option {
	return func(o *options) {
		o.EastAsian = enabled
	}
}

// WithSkipStrictEmojiNeutral should be set to true for some broken East-Asian fonts.
//
// This option only takes effect when EastAsianWidth is enabled.
func WithSkipStrictEmojiNeutral(enabled bool) Option {
	return func(o *options) {
		o.SkipStrictEmojiNeutral = enabled
	}
}

// WithAsianAmbiousWidth sets the default width to apply in EastAsian mode
// for ambiguous runes, i.e. non-EastAsian characters that already have a code point
// in an East-Asian character set.
//
// This option only takes effect when EastAsianWidth is enabled.
//
// The default is to apply 2 cells for those characters.
func WithAsianAmbiguousWidth(width int) Option {
	return func(o *options) {
		o.DefaultAsianAmbiguousWidth = width
	}
}

func optionsWithDefaults(opts []Option) *options {
	if len(opts) == 0 {
		return defaultOptions
	}

	o := *defaultOptions

	for _, apply := range opts {
		apply(&o)
	}

	return &o
}
