package ansi

import (
	"io"
	"unicode"

	"github.com/fredbi/go-typeset/terminal/runes"
)

const (
	esc     = '\033'
	bracket = '['
)

// StripANSIFromRunes strips a starting and a ending ANSI escape sequences in a slice of runes.
//
// If the slice of runes contains several sequences, the remaining runes after the end of the first end sequence are returned.
//
// This means that you may need to call StripANSIFromRunes several times until no start sequence is found to ensure that the stripped
// result no longer contains any sequence.
//
// It returns results in the following order:
// * the stripped string
// * the starting escape sequence(s)
// * the ending escape sequence(s)
// * reminder
func StripANSIFromRunes(rns []rune) ([]rune, []rune, []rune, []rune) {
	reader := runes.NewSliceReader(rns)
	var groups [3][]rune

	for {
		startANSI := decodeANSISequence(reader)
		if startANSI == nil {
			break
		}

		if len(groups[0]) == 0 {
			groups[0] = startANSI
		} else {
			// there are several sequences: indulge in 1 allocation
			groups[0] = append(groups[0], startANSI...)
		}
	}

	if str := decodeString(reader); str != nil {
		groups[1] = str
	}

	for {
		endANSI := decodeANSISequence(reader)
		if endANSI == nil {
			break
		}

		if len(groups[2]) == 0 {
			groups[2] = endANSI
		} else {
			groups[2] = append(groups[2], endANSI...)
		}
	}

	return groups[1], groups[0], groups[2], reader.Runes()
}

// decodeANSISequence identifies an ANSI terminal escape sequence.
//
// The escape sequence detected follows a Control Sequence Introducer: ESC[.
func decodeANSISequence(rdr *runes.SliceReader) []rune {
	const maxDigits = 2

	var (
		requiredMatch bool
		runeCount     int
		nextStart     int64
	)

	current, _ := rdr.Seek(0, io.SeekCurrent)

	defer func() {
		if !requiredMatch {
			_, _ = rdr.Seek(current, io.SeekStart)

			return
		}

		if nextStart > 0 {
			_, _ = rdr.Seek(nextStart, io.SeekStart)
		}
	}()

	r, _, err := rdr.ReadRune()
	if err != nil {
		return nil
	}
	if r != esc {
		return nil
	}
	runeCount++

	r, _, err = rdr.ReadRune()
	if err != nil {
		return nil
	}
	if r != bracket {
		return nil
	}
	runeCount++

	r, _, err = rdr.ReadRune()
	if err != nil {
		return nil
	}
	if !unicode.IsDigit(r) { // TODO: standard sayz is optional, as in ESC[m
		return nil
	}

	// required pattern found
	// Now matches the optional parts of the pattern
	requiredMatch = true
	runeCount++
	nextStart, _ = rdr.Seek(0, io.SeekCurrent)

	for i := 0; i < maxDigits; i++ {
		r, _, err = rdr.ReadRune()
		if err != nil {
			break
		}

		if !unicode.IsDigit(r) {
			break
		}

		runeCount++
		nextStart++
	}

	if r == ';' || r == ':' { // TODO: can have more that 2 sections, as in ESC[38,5;(n)m (256 bit color) or ESC[38;2;(r);(g);(b)m (Gnome term, kconsole)
		// TODO: ';' may be ':' in some case???
		runeCount++
		nextStart++

		r, _, err = rdr.ReadRune()
		if err != nil {
			return rdr.Slice(int(current), int(current)+runeCount)
		}

		if !unicode.IsDigit(r) { // TODO: is optional
			return rdr.Slice(int(current), int(current)+runeCount)
		}

		runeCount++
		nextStart++

		for i := 0; i < maxDigits; i++ {
			r, _, err = rdr.ReadRune()
			if err != nil {
				break
			}

			if !unicode.IsDigit(r) {
				break
			}

			runeCount++
			nextStart++
		}
	}

	// add the trailing keycode (e.g. m, K, ~...)
	runeCount++
	nextStart++

	return rdr.Slice(int(current), int(current)+runeCount)
}

func decodeString(rdr *runes.SliceReader) []rune {
	var (
		runeCount int
	)

	current, _ := rdr.Seek(0, io.SeekCurrent)

	for {
		r, _, err := rdr.ReadRune()
		if err != nil {
			break
		}

		if r == esc {
			start := current
			current, _ = rdr.Seek(-1, io.SeekCurrent)

			return rdr.Slice(int(start), int(start)+runeCount)
		}

		runeCount++
	}

	return rdr.Slice(int(current), int(current)+runeCount)
}

// StripANSI is like StriPANSIFromRunes but works with strings.
func StripANSI(str string) (string, string, string, string) {
	stripped, start, end, remainder := StripANSIFromRunes([]rune(str))

	return string(stripped), string(start), string(end), string(remainder)
}
