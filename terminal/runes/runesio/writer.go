package runesio

import (
	"io"
	"unicode/utf8"
)

type (
	// Writer knows how to write runes.
	Writer interface {
		WriteRunes([]rune) (int, error)
		WriteRune(rune) (int, error)
		io.Writer
	}

	runesWriter struct {
		io.Writer
	}
)

var _ Writer = &runesWriter{}

// NewWriter builds a new runes Writer from an io.Writer.
func NewWriter(w io.Writer) Writer {
	return &runesWriter{
		Writer: w,
	}
}

func (rw *runesWriter) WriteRune(r rune) (int, error) {
	return rw.writeRune(r)
}

func (rw *runesWriter) WriteRunes(text []rune) (int, error) {
	return rw.writeRunes(text)
}

func (rw *runesWriter) writeRune(r rune) (int, error) {
	// from strings.Builder.WriteRune
	var (
		buf [utf8.UTFMax]byte
		n   int
	)

	if uint32(r) < utf8.RuneSelf {
		buf[0] = byte(r)
		n = 1
	} else {
		n = utf8.EncodeRune(buf[:], r)
	}

	return rw.Writer.Write(buf[:n])
}

func (rw *runesWriter) writeRunes(runes []rune) (int, error) {
	buf := make([]byte, 0, utf8.UTFMax*len(runes))

	for _, r := range runes {
		if uint32(r) < utf8.RuneSelf {
			buf = append(buf, byte(r))

			continue
		}

		var bytesForRune [utf8.UTFMax]byte
		n := utf8.EncodeRune(bytesForRune[:], r)
		buf = append(buf, bytesForRune[:n]...)
	}

	return rw.Writer.Write(buf)
}
