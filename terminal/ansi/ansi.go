package ansi

import (
	"io"
	"unicode"

	"github.com/fredbi/go-typeset/terminal/runes/runesio"
)

const (
	esc            = '\033'
	openingBracket = '['
)

// StrippedToken represents a token stripped from ANSI escape sequences.
type StrippedToken struct {
	Text          []rune
	StartSequence []rune
	StopSequence  []rune
	Remainder     []rune
}

type kind uint8

const (
	otherSequence kind = iota
	startSequence
	stopSequence
)

// StripToken decodes a token into stripped tokens.
func StripToken(token []rune) []StrippedToken {
	attrs := make([]StrippedToken, 0, 4)

	for {
		stripped := StripANSIFromRunes(token)
		attrs = append(attrs, stripped)
		if len(stripped.Remainder) == 0 {
			break
		}

		token = stripped.Remainder
	}

	return attrs
}

// StripANSIFromRunes strips a starting and a ending ANSI escape sequences from a token provided as a slice of runes.
//
// If the slice of runes contains several sequences, the remaining runes after the end of the first end sequence are returned.
//
// This means that you may need to call StripANSIFromRunes several times until no start sequence is found to ensure that the stripped
// result no longer contains any sequence.
//
// Specifically, ANSI sequences detected as start/stop are SGR codes ("Select Graphic Rendition"):
//
//	ESC[{parameters}m
//
// Plus:
//
//	ESC[{parameters}s | ESC[{parameters}u
//	ESC[{parameters}h | ESC[{parameters}l
//
// It returns results in the following order:
// * the stripped string (possibly with some escape sequences that are neither start nor stop)
// * a detected starting escape sequence(s)
// * a detected ending escape sequence(s)
// * remainder
//
// TODO: resolve ambiguities such as when using "default color"
func StripANSIFromRunes(rns []rune) StrippedToken {
	reader := runesio.NewSliceReader(rns)
	var groups [3][]rune

	// decode starting sequence
	for {
		startANSI, sequenceKind := decodeANSISequence(reader)
		if startANSI == nil {
			break
		}

		// break when several sequences are chained or come out of natural order
		switch {
		case sequenceKind == otherSequence:
			// sequence not identified as start or stop: break it down separately
			return StrippedToken{
				Text:          startANSI,
				StartSequence: groups[0], // for when we have Start/Other in one single series of runes
				Remainder:     reader.Runes(),
			}

		case sequenceKind == startSequence && len(groups[0]) == 0:
			groups[0] = startANSI

		case sequenceKind == stopSequence:
			// token begins with a stop sequence: break in several results
			return StrippedToken{
				StartSequence: groups[0], // for when we have Start/Stop in one single series of runes
				StopSequence:  startANSI,
				Remainder:     reader.Runes(),
			}

		default:
			// there are several sequences: break in several results
			if len(startANSI) > 0 {
				// rewind on this sequence
				_, _ = reader.Seek(int64(-len(startANSI)), io.SeekCurrent)
			}

			return StrippedToken{
				StartSequence: groups[0],
				Remainder:     reader.Runes(),
			}
		}
	}

	// decode regular text (no esapce)
	if str := decodeText(reader); str != nil {
		groups[1] = str
	}

	// decode trailing sequence
	for {
		endANSI, sequenceKind := decodeANSISequence(reader)
		if endANSI == nil {
			break
		}

		// break when several sequences are chained or come out of natural order
		switch {
		case sequenceKind == otherSequence && len(groups[2]) == 0:
			// no Stop scanned yet
			if groups[1] == nil {
				groups[1] = endANSI // should normally never end up here
			} else {
				// indulge into one allocation
				groups[1] = append(groups[1], endANSI...)
			}

		case sequenceKind == otherSequence && len(groups[2]) > 0:
			// comes after a Stop
			// rewind on this sequence
			_, _ = reader.Seek(int64(-len(endANSI)), io.SeekCurrent)

			return StrippedToken{
				Text:          groups[1],
				StartSequence: groups[0],
				StopSequence:  groups[2],
				Remainder:     reader.Runes(),
			}

		case sequenceKind == stopSequence && len(groups[2]) == 0:
			groups[2] = endANSI

		case sequenceKind == startSequence:
			// token ends with a start sequence: break in several results
			// rewind on this sequence
			_, _ = reader.Seek(int64(-len(endANSI)), io.SeekCurrent)

			return StrippedToken{
				Text:          groups[1],
				StartSequence: groups[0],
				StopSequence:  groups[2],
				Remainder:     reader.Runes(),
			}

		default:
			// there are several sequences: break in several results
			if len(endANSI) > 0 {
				// rewind on this sequence
				_, _ = reader.Seek(int64(-len(endANSI)), io.SeekCurrent)
			}

			return StrippedToken{
				Text:          groups[1],
				StartSequence: groups[0],
				StopSequence:  groups[2],
				Remainder:     reader.Runes(),
			}
		}
	}

	return StrippedToken{
		Text:          groups[1],
		StartSequence: groups[0],
		StopSequence:  groups[2],
		Remainder:     reader.Runes(),
	}
}

// decodeANSISequence identifies an ANSI terminal escape sequence.
//
// The escape sequence detected follows a Control Sequence Introducer: ESC[.
func decodeANSISequence(rdr *runesio.SliceReader) ([]rune, kind) {
	const maxDigits = 3

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

	// require escape sequence start
	r, _, err := rdr.ReadRune()
	if err != nil {
		return nil, otherSequence
	}
	if r != esc {
		return nil, otherSequence
	}
	runeCount++

	r, _, err = rdr.ReadRune()
	if err != nil {
		return nil, otherSequence
	}
	if r != openingBracket {
		return nil, otherSequence
	}
	runeCount++

	// required pattern found
	// Now matches the optional parts of the pattern
	requiredMatch = true
	nextStart, _ = rdr.Seek(0, io.SeekCurrent)

	hasZeroArg := true
	for i := 0; i < maxDigits; i++ {
		r, _, err = rdr.ReadRune()
		if err != nil {
			break
		}

		if !unicode.IsDigit(r) {
			break
		}

		if r != '0' {
			hasZeroArg = false
		}

		runeCount++
		nextStart++
	}

	for r == ';' || r == ':' { // can have more that 2 sections, as in ESC[38,5;(n)m (256 bit color) or ESC[38;2;(r);(g);(b)m (Gnome term, kconsole)
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
	var sequenceKind kind
	switch r {
	case 's', 'h':
		sequenceKind = startSequence
	case 'u', 'l':
		sequenceKind = stopSequence
	case 'm':
		// standard SGR
		if hasZeroArg {
			sequenceKind = stopSequence
		} else {
			sequenceKind = startSequence
		}
	default:
		sequenceKind = otherSequence
	}

	if err == nil {
		runeCount++
		nextStart++
	}

	return rdr.Slice(int(current), int(current)+runeCount), sequenceKind
}

func decodeText(rdr *runesio.SliceReader) []rune {
	var runeCount int

	current, _ := rdr.Seek(0, io.SeekCurrent)

	for {
		r, _, err := rdr.ReadRune()
		if err != nil {
			break
		}

		if r == esc {
			r, _, err = rdr.ReadRune()
			if err != nil {
				runeCount++

				break
			}

			if r == openingBracket {
				_, _ = rdr.Seek(-2, io.SeekCurrent)

				// start of another escape sequence: end of stripped string
				return rdr.Slice(int(current), int(current)+runeCount)
			}

			runeCount++
		}

		runeCount++
	}

	return rdr.Slice(int(current), int(current)+runeCount)
}
