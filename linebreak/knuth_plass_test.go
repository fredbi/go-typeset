package linebreak

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fredbi/go-typeset/wordbreak/hyphenator"
	"github.com/stretchr/testify/require"
)

func TestLineBreaker(t *testing.T) {
	t.Run("with line length matching text, shoud left-align", func(t *testing.T) {
		const (
			paragraph = `Lorem ipsum dolor sit amet, consectetur adipiscing` // matches exactly 1 line
			display   = 50.0
		)
		t.Run("render", testLeftAlign(paragraph, display))
	})

	t.Run("with lorem ipsum", func(t *testing.T) {
		const paragraph = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, ` +
			`sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. ` +
			`Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. ` +
			`Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. ` +
			`Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

		t.Run("should left-align a paragraph (width 50)", func(t *testing.T) {
			const display = 50.0
			t.Run("render", testLeftAlign(paragraph, display))

			t.Run("render with hyphens (width 40)", testLeftAlignHyphenize(paragraph, display-10))
		})
		t.Run("should left-align a paragraph (width 51)", func(t *testing.T) {
			const display = 51.0
			t.Run("render", testLeftAlign(paragraph, display))

			t.Run("render with hyphens (width 41)", testLeftAlignHyphenize(paragraph, display-10))
		})
	})

	t.Run("should left-align Knuth's classical example (width 25)", func(t *testing.T) {
		const (
			paragraph = `In olden times when wishing still helped one, there lived a king ` +
				`whose daughters were all beautiful, but the youngest was so beautiful ` +
				`that the sun itself, which has seen so much, was astonished whenever it ` +
				`shone in her face. Close by the king's castle lay a great dark forest, ` +
				`and under an old lime-tree in the forest was a well, and when the day ` +
				`was very warm, the king's child went out into the forest and sat down by ` +
				`the side of the cool fountain, and when she was bored she took a golden ball, ` +
				`and threw it up on high and caught it, and this ball was her favorite plaything.`
			display = 30.0
		)

		t.Run("render", testLeftAlign(paragraph, display))
		/// TODO: assert Knuth's breakpoints

		t.Run("render with hyphens (width 20)", testLeftAlignHyphenize(paragraph, display-10))
	})

	t.Run("should hyphenate long text (width 8)", func(t *testing.T) {
		const (
			paragraph = `In olden times when wishing still helped one, there lived a king ` +
				`whose daughters were all beautiful, but the youngest was so beautiful ` +
				`that the sun itself, which has seen so much, was astonished whenever it ` +
				`shone in her face. Close by the king's castle lay a great dark forest, ` +
				`and under an old lime-tree in the forest was a well, and when the day ` +
				`was very warm, the king's child went out into the forest and sat down by ` +
				`the side of the cool fountain, and when she was bored she took a golden ball, ` +
				`and threw it up on high and caught it, and this ball was her favorite plaything.`
			display = 8.0
		)

		t.Run("render with hyphens", testLeftAlignHyphenize(paragraph, display))
	})

	t.Run("should hyphenate a super long word (width 20)", func(t *testing.T) {
		const (
			//           0123456789012345678901234567890123456789
			//                              |
			paragraph = `Honorificabilitudinitatibus is super long`
			/*
				            |Honorificabilitudi-.|
				            |nitatibus is super..|
							|long
			*/
			display = 20.0
		)

		t.Run("render with hyphens", testLeftAlignHyphenize(paragraph, display))
	})
}

func testLeftAlign(paragraph string, display float64) func(*testing.T) {
	return func(t *testing.T) {
		tokens := strings.Fields(paragraph)

		lb := New() // with all defaults
		lines, err := lb.LeftAlignUniform(tokens, display)
		require.NoError(t, err)

		testRenderLines(lines, display)
	}
}

func testRenderLines(lines []string, display float64) {
	for _, line := range lines {
		pad := strings.Repeat(".", max(0, int(display)-len(line)))
		fmt.Printf("|%s%s|\n", line, pad)
	}
}

func testLeftAlignHyphenize(paragraph string, display float64) func(*testing.T) {
	return func(t *testing.T) {
		tokens := strings.Fields(paragraph)

		h := hyphenator.New()
		lb := New(
			WithRenderHyphens(true),
			WithHyphenator(h.BreakWord),
			WithHyphenPenalty(50.00),
		)
		lines, err := lb.LeftAlignUniform(tokens, display)
		require.NoError(t, err)

		testRenderLines(lines, display)
	}
}
