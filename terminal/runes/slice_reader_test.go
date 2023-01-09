package runes

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSliceReader(t *testing.T) {
	t.Parallel()
	const str = "longstring"

	t.Run("should iterate over a set of runes", func(t *testing.T) {
		t.Parallel()

		r := NewSliceReader([]rune(str))

		var i int
		res := make([]rune, 0, len(str))

		for {
			rn, size, err := r.ReadRune()
			if err != nil {
				break
			}
			res = append(res, rn)
			i++

			require.Equal(t, Width(rn), size)
		}

		require.Equal(t, len(str), i)
		require.Equal(t, str, string(res))
	})

	t.Run("should Seek over a set of runes", func(t *testing.T) {
		r := NewSliceReader([]rune(str))

		t.Run("should Seek to end", func(t *testing.T) {
			offset, err := r.Seek(int64(len(str)), io.SeekStart)
			require.NoError(t, err)
			require.Equal(t, len(str), int(offset))

			_, _, err = r.ReadRune()
			require.ErrorIs(t, err, io.EOF)
		})

		t.Run("with different whence starting points", func(t *testing.T) {
			t.Parallel()

			t.Run("should Seek from Start", func(t *testing.T) {
				offset, err := r.Seek(1, io.SeekStart)
				require.NoError(t, err)
				require.Equal(t, 1, int(offset))

				rn, _, err := r.ReadRune()
				require.NoError(t, err)
				const expected = 'o'
				require.Equalf(t, expected, rn,
					"expected %c but got %c",
					expected, rn,
				)
			})

			t.Run("should Seek from current", func(t *testing.T) {
				rn, _, err := r.ReadRune()
				require.NoError(t, err)
				const expected = 'n'
				require.Equalf(t, expected, rn,
					"expected %c but got %c",
					expected, rn,
				)

				offset, err := r.Seek(1, io.SeekCurrent)
				require.NoError(t, err)
				require.Equal(t, 4, int(offset))

				rn, _, err = r.ReadRune()
				require.NoError(t, err)
				const expectedNext = 's'
				require.Equalf(t, expectedNext, rn,
					"expected %c but got %c",
					expectedNext, rn,
				)
			})

			t.Run("should Seek from end", func(t *testing.T) {
				offset, err := r.Seek(-1, io.SeekEnd)
				require.NoError(t, err)
				require.Equal(t, len(str)-1, int(offset))

				rn, _, err := r.ReadRune()
				require.NoError(t, err)
				const expected = 'g'
				require.Equalf(t, expected, rn,
					"expected %c but got %c",
					expected, rn,
				)
			})
		})
	})

	t.Run("Seek should error on out of range offsets", func(t *testing.T) {
		t.Parallel()

		r := NewSliceReader([]rune(str))

		_, err := r.Seek(-1, io.SeekStart)
		require.Error(t, err)

		_, err = r.Seek(int64(len(str)+1), io.SeekStart)
		require.Error(t, err)

		_, err = r.Seek(int64(len(str)+1), io.SeekCurrent)
		require.Error(t, err)
	})
}
