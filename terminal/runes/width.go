//go:generate go run build_tables.go

package runes

import "unicode/utf8"

// Width returns the number of cells in a single rune.
//
// NOTE: this version does not support graphemes on multiple code points.
//
// References:
//   - East-Asian characters are displayed as: http://www.unicode.org/reports/tr11
//   - Emojis are displayed as: https://www.unicode.org/reports/tr51
func Width(r rune, opts ...Option) int {
	if !utf8.ValidRune(r) || r == utf8.RuneError {
		return 0
	}

	o := optionsWithDefaults(opts)
	if o.EastAsian || o.SkipStrictEmojiNeutral {
		return runeWidth(r, o)
	}

	buildLookupOnce.Do(initLookupTable)

	return int(
		(lookupTable[r>>2] >> ((r % 4) * 2)) & 3,
	)
}

// Widths return the width of a slice of runes as displayed in fixed-width font.
//
// NOTE: this version does not support graphemes on multiple code points.
func Widths(runes []rune, opts ...Option) (width int) {
	o := optionsWithDefaults(opts)

	for _, r := range runes {
		width += runeWidth(r, o)
	}

	return width
}

// StringWidth return the width of a string, as displayed in fixed-width font.
func StringWidth(s string, opts ...Option) (width int) {
	return Widths([]rune(s), opts...)
}

func runeWidth(r rune, opts *options) int {
	if !utf8.ValidRune(r) || r == utf8.RuneError {
		return 0
	}

	if opts.EastAsian {
		return runeWidthEastAsian(r, opts.SkipStrictEmojiNeutral)
	}

	return runeWidthNonEastAsian(r)
}

func runeWidthNonEastAsian(r rune) int {
	switch {
	case r < 0x20:
		return 0
	case (r >= 0x7F && r <= 0x9F) || r == 0xAD: // non-printable
		return 0
	case r < 0x300:
		return 1
	case inTables(r, nonprint, combining):
		return 0
	case inTable(r, narrow):
		return 1
	case inTable(r, doublewidth):
		return 2
	default:
		return 1
	}
}

func runeWidthEastAsian(r rune, skipStrictEmojiNeutral bool) int {
	switch {
	case inTables(r, nonprint, combining):
		return 0
	case inTable(r, narrow):
		return 1
	case inTables(r, ambiguous, doublewidth):
		return 2
	case skipStrictEmojiNeutral && inTables(r, ambiguous, emoji, narrow):
		// TODO: case skipStrictEmojiNeutral && inTables(r, emoji, narrow):
		return 2
	default:
		return 1
	}
}
