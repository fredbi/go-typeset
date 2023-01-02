package separator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeparator(t *testing.T) {
	t.Run("should break paths", func(t *testing.T) {
		const word = `a/b/c|d-e_f\g-`

		s := New()
		require.Equal(t, []string{
			"a/", "b/", "c|", "d-", "e_", `f\`, "g-",
		},
			s.BreakWord(word),
		)
	})
}
