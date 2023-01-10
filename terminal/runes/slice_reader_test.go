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

func TestSlicing(t *testing.T) {
	const str = "unixöperatingsystem"

	t.Run("should slice the collection of runes with byte offsets", func(t *testing.T) {
		r := NewSliceReader([]rune(str))
		slices := r.SliceFromByteOffsets([]int{1, 4})
		expected := []string{"nix"}
		require.Equal(t, expected, toStrings(slices))

		slices = r.SliceFromByteOffsets([]int{4, 8, 10, 12})
		expected = []string{"öpe", "ti"}
		require.Equal(t, expected, toStrings(slices))
	})

	t.Run("should slice the collection of runes with rune offsets", func(t *testing.T) {
		r := NewSliceReader([]rune(str))
		slices := r.Slices([]int{1, 4})
		expected := []string{"nix"}
		require.Equal(t, expected, toStrings(slices))

		slices = r.Slices([]int{4, 8, 10, 12})
		expected = []string{"öper", "in"}
		require.Equal(t, expected, toStrings(slices))
	})

	t.Run("should slice the collection of runes with single rune offset", func(t *testing.T) {
		r := NewSliceReader([]rune(str))
		slice := r.Slice(1, 4)
		expected := "nix"
		require.Equal(t, expected, string(slice))

		slice = r.Slice(4, 8)
		expected = "öper"
		require.Equal(t, expected, string(slice))
	})

	t.Run("should panic if offsets are not pairs", func(t *testing.T) {
		r := NewSliceReader([]rune(str))

		require.Panics(t, func() {
			_ = r.SliceFromByteOffsets([]int{4, 10, 12}) //nolint:staticcheck
		})

		require.Panics(t, func() {
			_ = r.Slices([]int{4, 8, 10}) //nolint:staticcheck
		})
	})

	t.Run("should keep track of unread runes so far", func(t *testing.T) {
		r := NewSliceReader([]rune(str))
		for i := 0; i < 5; i++ {
			_, _, _ = r.ReadRune()
		}
		require.Equal(t, "peratingsystem", string(r.Runes()))
	})
}

func toStrings(in [][]rune) []string {
	out := make([]string, 0, len(in))

	for _, s := range in {
		out = append(out, string(s))
	}

	return out
}
