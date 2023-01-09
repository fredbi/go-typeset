package runes

import (
	"errors"
	"io"
	"unicode/utf8"
)

var (
	_ io.RuneReader = &SliceReader{}
)

type SliceReader struct {
	offset int
	runes  []rune
}

func NewSliceReader(runes []rune) *SliceReader {
	return &SliceReader{runes: runes}
}

func (r *SliceReader) ReadRune() (rune, int, error) {
	if r.offset >= len(r.runes) {
		return utf8.RuneError, 0, io.EOF
	}

	rn := r.runes[r.offset]
	r.offset++

	return rn, utf8.RuneLen(rn), nil
}

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
		return -1, errors.New("Seek: out of range offset")
	}

	r.offset = int(pos)

	return pos, nil
}
