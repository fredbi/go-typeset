package runesio

import (
	"io"
	"unicode/utf8"
)

type (
	Writer interface {
		WriteRunes([]rune) (int, error)
		WriteRune(rune) (int, error)
		io.Writer
	}

	runesWriter struct {
		io.Writer
		writeRune  func(rune) (int, error)
		writeRunes func([]rune) (int, error)
	}
)

var _ Writer = &runesWriter{}

func AsRunesWriter(w io.Writer) Writer {
	return &runesWriter{
		Writer: w,
		// from strings.Builder.WriteRune
		writeRune: func(r rune) (int, error) {
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

			return w.Write(buf[:n])
		},
		writeRunes: func(runes []rune) (int, error) {
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

			return w.Write(buf)
		},
	}
}

func (rw *runesWriter) WriteRune(r rune) (int, error) {
	return rw.writeRune(r)
}

func (rw *runesWriter) WriteRunes(text []rune) (int, error) {
	return rw.writeRunes(text)
}
