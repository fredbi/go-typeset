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

// WithEastAsian should be set to true when displaying Asian characters.
//
// Notice that wide EastAsian runes are always corrected reported, with or without this option.
//
// The difference is only for "ambiguous" runes that appear in different character sets.
// Those normally reported with width 1 are reported with width 2 with this option enabled.
//
// Example: 'ø' is reported with width 1 by default (EastAsianWidth disabled)
// and with width 2 when EastAsianWidth is enabled.
//
// You may change the default width of ambiguous runes with the WithEastAsianAmbiguousWith() option.
//
// NOTE: the current implementation is optimized with this mode disabled. A faster lookup is performed.
func WithEastAsian(enabled bool) Option {
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

// WithEastAsianAmbiousWidth sets the default width to apply in EastAsian mode
// for ambiguous runes, i.e. non-EastAsian characters that already have a code point
// in an East-Asian character set.
//
// This option only takes effect when EastAsianWidth is enabled.
//
// The default is to apply 2 cells for those characters.
func WithEastAsianAmbiguousWidth(width int) Option {
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
