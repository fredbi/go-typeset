//nolint:forbidigo
package linebreak

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fredbi/go-typeset/terminal/ansi"
	"github.com/fredbi/go-typeset/terminal/runes"
	"github.com/fredbi/go-typeset/wordbreak/hyphenator"
	"github.com/stretchr/testify/require"
)

// This excerpt from The Frog King (Grimm) pays tribute to the long-established tradition
// to include Knuth's original test cases in all variations of this algorithm.
const grimm = `In olden times when wishing still helped one, there lived a king ` +
	`whose daughters were all beautiful, but the youngest was so beautiful ` +
	`that the sun itself, which has seen so much, was astonished whenever it ` +
	`shone in her face. Close by the king's castle lay a great dark forest, ` +
	`and under an old lime-tree in the forest was a well, and when the day ` +
	`was very warm, the king's child went out into the forest and sat down by ` +
	`the side of the cool fountain, and when she was bored she took a golden ball, ` +
	`and threw it up on high and caught it, and this ball was her favorite plaything.`

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
			paragraph = grimm
			display   = 30.0
		)

		t.Run("render", testLeftAlign(paragraph, display))
		/// TODO: assert Knuth's breakpoints

		t.Run("render with hyphens (width 20)", testLeftAlignHyphenize(paragraph, display-10))
	})

	t.Run("should hyphenate long text (width 8)", func(t *testing.T) {
		const (
			paragraph = grimm
			display   = 8.0
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

	t.Run("should hyphenate long text (width 12)", func(t *testing.T) {
		const (
			paragraph = grimm
			display   = 12.0
		)

		t.Run("render with hyphens", testLeftAlignHyphenize(paragraph, display))
	})

	t.Run("should hyphenate a super long word (width 15)", func(t *testing.T) {
		const (
			//           0123456789012345678901234567890123456789
			//                              |
			paragraph = `Honorificabilitudinitatibus is super long`
			// |Honorificabil-.|
			// |itudinitatibus.|
			// |is super long..|
			//
			display = 15.0
		)

		t.Run("render with hyphens", testLeftAlignHyphenize(paragraph, display))
	})

	t.Run("should detect separators and split accordingly", func(t *testing.T) {
		const (
			paragraph = `https://www.unicode.org/Public/15.0.0/ucd/emoji/emoji-data.txt`
			display   = 15.0
		)

		t.Run("render with hyphens", testLeftAlignHyphenize(paragraph, display, WithRenderHyphens(false)))
	})

	t.Run("with attributes", func(t *testing.T) {
		const (
			startRed   = "\033[31;1;4m"
			stop       = "\033[0m"
			startGreen = "\033[92m"

			paragraph = `Lorem ipsum dolor ` + startRed + `sit amet, consectetur adipiscing elit, ` +
				`sed do eiusmod tempor incididunt ut labore et dolore` + stop + ` magna aliqua. ` +
				`Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. ` +
				`Duis aute irure dolor in reprehenderit ` + startGreen + `in voluptate velit esse cillum dolore eu fugiat nulla pariatur. ` +
				`Excepteur sint occaecat cupidatat non proident,` + stop + ` sunt in culpa qui officia deserunt mollit anim id est laborum.`
		)

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

	t.Run("with nested attributes", func(t *testing.T) {
		t.SkipNow() // not supported for now
		const (
			startUnderline = "\033[4m"
			stopUnderline  = "\033[24m"
			startRed       = "\033[91m"
			stopRed        = "\033[39m"

			paragraph = `Lorem ipsum dolor ` + startRed + `sit amet, consectetur ` + startUnderline + `adipiscing elit, ` + stopUnderline +
				`sed do eiusmod tempor incididunt ut labore et dolore` + stopRed + ` magna aliqua. ` +
				`Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. ` +
				`Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. ` +
				`Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`
		)

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
}

//nolint:unparam
func testLeftAlign(paragraph string, display float64, opts ...Option) func(*testing.T) {
	return func(t *testing.T) {
		tokens := strings.Fields(paragraph)
		opts = append([]Option{WithWordBreak(false)}, opts...)

		lb := New(opts...)
		lines, err := lb.LeftAlignUniform(tokens, display)
		require.NoError(t, err)

		testRenderLines(lines, display)
	}
}

func testRenderLines(lines []string, display float64) {
	for _, line := range lines {
		var w int
		for _, stripped := range ansi.StripToken([]rune(line)) {
			w += runes.Widths(stripped.Text)
		}
		pad := strings.Repeat("*", max(0, int(display)-w))
		fmt.Printf("|%s%s|\n", line, pad)
	}
}

func testLeftAlignHyphenize(paragraph string, display float64, opts ...Option) func(*testing.T) {
	return func(t *testing.T) {
		tokens := strings.Fields(paragraph)
		h := hyphenator.New()
		opts = append([]Option{
			WithRenderHyphens(true),
			WithHyphenator(h.BreakWord),
		}, opts...)

		lb := New(opts...)
		lines, err := lb.LeftAlignUniform(tokens, display)
		require.NoError(t, err)

		testRenderLines(lines, display)
	}
}
