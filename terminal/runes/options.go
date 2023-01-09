package runes

type (
	// Option to set expectations on rune width computation behavior.
	Option func(*options)

	options struct {
		EastAsian              bool
		SkipStrictEmojiNeutral bool
	}
)

var (
	defaultOptions = &options{
		EastAsian:              false,
		SkipStrictEmojiNeutral: false,
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

// WithSkipStrictEmojiNeutral should be set to true for some broken East-Asian fonts
func WithSkipStrictEmojiNeutral(enabled bool) Option {
	return func(o *options) {
		o.SkipStrictEmojiNeutral = enabled
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
