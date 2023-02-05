package runesio

import (
	"errors"
	"io"
	"unicode/utf8"
)

var (
	_ io.RuneReader = &SliceReader{}
	_ io.Seeker     = &SliceReader{}

	// ErrSeekOutOfRange complains that an attempt was made to Seek out of range.
	ErrSeekOutOfRange = errors.New("Seek: out of range offset")
)

// SliceReader is an io.RuneReader from a slice of runes.
type SliceReader struct {
	offset int
	runes  []rune
}

// NewSliceReader builds a SliceReader from a slice of runes.
func NewSliceReader(runes []rune) *SliceReader {
	return &SliceReader{runes: runes}
}

// ReadRune consumes a single rune.
//
// The size read is expressed in bytes.
func (r *SliceReader) ReadRune() (rune, int, error) {
	if r.offset >= len(r.runes) {
		return utf8.RuneError, 0, io.EOF
	}

	rn := r.runes[r.offset]
	r.offset++

	return rn, utf8.RuneLen(rn), nil
}

// Runes returns unread runes.
func (r *SliceReader) Runes() []rune {
	return r.runes[r.offset:]
}

// Seek positions the reader at some offset. Offset is expressed in runes, not in bytes.
//
// The "whence" argument works as defined by io.Seek and is one of io.SeekStart, io.SeekdEnd, io.SeekCurrent.
func (r *SliceReader) Seek(offset int64, whence int) (int64, error) {
	pos := int64(r.offset)
	switch whence {
	case io.SeekStart:
		pos = offset
	case io.SeekEnd:
		pos = int64(len(r.runes)) + offset
	case io.SeekCurrent:
		pos += offset
	}

	if pos < 0 || pos > int64(len(r.runes)) {
		return -1, ErrSeekOutOfRange
	}

	r.offset = int(pos)

	return pos, nil
}

// Slice returns the underlying slice of runes.
//
// It operates from any state of the reader.
func (r *SliceReader) Slice(startRune, endRune int) []rune {
	return r.runes[startRune:endRune]
}

// Slices returns a collection of slices of the underlying runes,
// with cut points specifyed by the pairs of offsets provided.
//
// Offsets are provided in runes.
//
// It operates from any state of the reader.
//
// Slices  panics if the slice argument does not come in pairs.
func (r *SliceReader) Slices(slice []int) [][]rune {
	if len(slice)%2 != 0 {
		panic("provided slices should come in pairs")
	}

	groups := make([][]rune, 0, len(slice)/2)

	for index := 0; index <= len(slice)-2; index += 2 {
		startRune := slice[index]
		endRune := slice[index+1]
		groups = append(groups, r.runes[startRune:endRune])
	}

	return groups
}

// SliceFromByteOffsets returns slices of the underlying set of runes,
// using pairs of offsets expressed in bytes, e.g. from a regexp
// operating on the RuneReader.
//
// It operates from any state of the reader.
//
// SliceFromByteOffsets panics if the slice argument does not come in pairs.
func (r *SliceReader) SliceFromByteOffsets(slice []int) [][]rune {
	if len(slice)%2 != 0 {
		panic("provided slices should come in pairs")
	}

	groups := make([][]rune, 0, len(slice)/2)

	for index := 0; index <= len(slice)-2; index += 2 {
		startByte := slice[index]
		endByte := slice[index+1]
		var startRune, endRune, offset int

		for _, rn := range r.runes {
			if offset < startByte {
				startRune++
			}
			offset += utf8.RuneLen(rn)

			if offset > endByte {
				break
			}

			endRune++
		}
		groups = append(groups, r.runes[startRune:endRune])
	}

	return groups
}
