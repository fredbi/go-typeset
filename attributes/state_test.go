package attributes

import (
	"fmt"
	"os"
	"testing"

	"github.com/fredbi/go-typeset/terminal/runes/runesio"
	"github.com/stretchr/testify/require"
)

func TestState(t *testing.T) {
	// TODO: assertions
	w := os.Stdout
	writer := runesio.AsRunesWriter(w)

	t.Run("with nested levels", func(t *testing.T) {
		state := NewState()
		testCases := []*Attribute{
			New([]rune("first"), []rune("\\033[start1"), nil),
			New([]rune("second"), nil, nil),
			New([]rune("third"), []rune("\\033[start2"), nil),
			New([]rune("fourth"), nil, []rune("\\033[stop2")),
			New([]rune("fifth"), nil, []rune("\\033[stop1")),
			New([]rune("sixth"), nil, nil),
		}

		for _, toPin := range testCases {
			attr := toPin
			state.Push(attr)
		}

		t.Run("check the inner chained list", func(t *testing.T) {
			require.Equal(t, len(testCases), state.list.Len())

			head := state.list.Front()
			require.NotNil(t, head)
			first, ok := head.Value.(Renderer)
			require.True(t, ok)
			require.Equal(t, "first", string(first.Runes()))
			require.Equal(t, 1, first.Level())

			next := head.Next()
			require.NotNil(t, next)
			second, ok := next.Value.(Renderer)
			require.True(t, ok)
			require.Equal(t, "second", string(second.Runes()))
			require.Equal(t, 1, second.Level())
		})

		t.Run("should iterate over the chained list", func(t *testing.T) {
			i := 0
			for iter := state.Iterator(); iter.Next(); i++ {
				item := iter.Item()
				item.Render(writer)

				if i == 3 { // after token 'fourth'
					iter.EndOfLine(writer)
					fmt.Fprintln(w, "")
					iter.StartOfLine(writer)
				}
			}
			require.Equal(t, len(testCases), i)
		})
	})

	t.Run("with start and stop in single attr", func(t *testing.T) {
		state := NewState()
		testCases := []*Attribute{
			New([]rune("first"), []rune("\\033[start1"), nil),
			New([]rune("second"), []rune("\\033[start3"), []rune("\\033[stop3")),
			New([]rune("third"), []rune("\\033[start2"), nil),
			New([]rune("fourth"), nil, []rune("\\033[stop2")),
			New([]rune("fifth"), nil, []rune("\\033[stop1")),
			New([]rune("sixth"), nil, nil),
		}

		for _, toPin := range testCases {
			attr := toPin
			state.Push(attr)
		}

		t.Run("should iterate over the chained list", func(t *testing.T) {
			i := 0
			for iter := state.Iterator(); iter.Next(); i++ {
				item := iter.Item()
				item.Render(writer)

				if i == 3 { // after token 'fourth'
					iter.EndOfLine(writer)
					fmt.Fprintln(w, "")
					iter.StartOfLine(writer)
				}
			}
			require.Equal(t, len(testCases), i)
		})
	})
}
