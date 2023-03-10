//go:generate go run build_tables.go

package runes

import (
	"unicode/utf8"
)

// Width returns the number of cells in a single rune.
//
// NOTE: this version does not support graphemes on multiple code points.
//
// References:
//   - East-Asian characters are displayed as: https://www.unicode.org/reports/tr11
//   - Emojis are displayed as: https://www.unicode.org/reports/tr51
func Width(r rune, opts ...Option) int {
	if !utf8.ValidRune(r) || r == utf8.RuneError {
		return 0
	}
	if r == 0x2E3B {
		// special case with code point with a width that overflow our lookup table
		return 4
	}

	o := optionsWithDefaults(opts)
	if o.EastAsian {
		return runeWidth(r, o)
	}

	buildLookupOnce.Do(initLookupTable) // builds the cache once

	return int(
		(lookupTable[r>>2] >> ((r % 4) * 2)) & 3,
	)
}

// Widths return the width of a slice of runes as displayed in fixed-width font.
//
// NOTE: this version does not support graphemes on multiple code points.
func Widths(runes []rune, opts ...Option) (width int) {
	o := optionsWithDefaults(opts)

	if o.EastAsian {
		for _, r := range runes {
			width += runeWidth(r, o)
		}

		return width
	}

	buildLookupOnce.Do(initLookupTable) // builds the cache once
	for _, r := range runes {
		width += int((lookupTable[r>>2] >> ((r % 4) * 2)) & 3)
	}

	return width
}

// StringWidth return the width of a string, as displayed in fixed-width font.
func StringWidth(s string, opts ...Option) (width int) {
	return Widths([]rune(s), opts...)
}

// IsAmbiguous returns true if the unicode point for this rune
// is considered ambiguous regarding width in East Asian character sets.
//
// Ambiguous runes refer to unicode points that are presents in different sets,
// and more context might be needed to determine the desired width.
func IsAmbiguous(r rune) bool {
	return inTable(r, ambiguous)
}

func runeWidth(r rune, opts *options) int {
	switch {
	case !utf8.ValidRune(r) || r == utf8.RuneError:
		return 0
	case r < 0x20:
		return 0
	case (r >= 0x7F && r <= 0x9F) || r == 0xAD: // non-printable
		return 0
	case r < 0x300 && !opts.EastAsian: // those code points are mostly ambiguous in EastAsian mode
		return 1
	case r == 0x2E3B: // THREE-EM-DASH
		return 4
	case r == 0x2E3A: // TWO-EM DASH
		return 3
	case inTables(r, nonprint, combining):
		return 0
	case inTable(r, narrow):
		return 1
	case inTable(r, doublewidth):
		return 2
	default:
		if !opts.EastAsian {
			return 1
		}

		// East-Asian behavior
		if inTables(r, ambiguous) {
			return opts.DefaultAsianAmbiguousWidth
		}

		if opts.SkipStrictEmojiNeutral && inTables(r, emoji, narrow) {
			return 2
		}

		return 1
	}
}
