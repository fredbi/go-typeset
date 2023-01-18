package linebreak

import (
	"container/list"
	"math"
	"strings"

	"github.com/fredbi/go-typeset/attributes"
	"github.com/fredbi/go-typeset/terminal/ansi"
	"github.com/fredbi/go-typeset/terminal/runes/runesio"
	"github.com/fredbi/go-typeset/wordbreak/hyphenator"
	"github.com/fredbi/go-typeset/wordbreak/punctuator"
)

type (
	// LineBreaker breaks a paragraph into lines under line width constraints.
	//
	// It implements the classical Knuth-Plass algorithm.
	LineBreaker struct {
		nodes       []nodeT
		lineWidths  []float64
		sum         *sums
		activeNodes *list.List
		spaceWidth  float64
		hyphenWidth float64
		attrList    *attributes.State
		// spaceStretch float64
		// spaceShrink  float64

		*options
	}
)

const (
	flaggedPenalty   = true
	unflaggedPenalty = false
	noWidth          = 0.0
	noShrink         = 0.0
)

var (
	space  = []rune{' '}
	hyphen = []rune{'-'} // rendering happens with a hard hyphen (visible)
)

// New line breaker.
func New(opts ...Option) *LineBreaker {
	l := &LineBreaker{
		attrList: attributes.NewState(),
		options:  defaultOptions(opts),
	}

	l.spaceWidth = l.scale(l.measurer(space))
	if l.wordBreak {
		l.hyphenWidth = l.scale(l.measurer(hyphen))
		if l.hyphenator == nil {
			// default hyphenator
			h := hyphenator.New(hyphenator.WithMinLength(l.minHyphenate))
			l.hyphenator = h.BreakWord
		}
	}

	if l.punctuator == nil {
		// default punctuator
		p := punctuator.New()
		l.punctuator = p.BreakWord
	}

	// NOTE: this is for center & justify (not implemented for now)
	// l.spaceStretch = min(1, int(float64(l.spaceWidth*l.space.width)/float64(l.space.stretch)))
	// l.spaceShrink = min(1, int(float64(l.spaceWidth*l.space.width)/float64(l.space.shrink)))

	return l
}

// LeftAlignUniform left-align a series of tokens that compose a paragraph,
// rendering multiple lines of uniform length maxLength.
//
// TODO: move to [][]rune -> []rune
func (l *LineBreaker) LeftAlignUniform(tokens []string, maxWidth float64) ([]string, error) {
	// 1. build a model that represent the tokens in terms of glue/box/penalty nodes
	l.nodes = l.leftAlignedNodes(tokens)

	// 2. build a model for desired widths for lines
	l.lineWidths = l.buildUniformLengths(maxWidth)

	// 3. compute a chained-list of break points
	breakList := l.breakPoints()
	if breakList == nil {
		return nil, ErrCannotBeSet
	}

	// 4. render the final wrapped text
	return l.render(breakList), nil
}

// buildUniformLengths fills 1 line of line length constraints,
// making this constraint uniform across all lines.
func (l *LineBreaker) buildUniformLengths(width float64) []float64 {
	lengths := make([]float64, 1)
	for i := range lengths {
		lengths[i] = l.scale(width)
	}

	return lengths
}

// scale a text measure to convert the measurement unit to a value compatible
// with the algorithm's parameters.
func (l *LineBreaker) scale(in float64) float64 {
	return math.Round(in * l.scaleFactor)
}

func (l *LineBreaker) downScale(in float64) float64 {
	return math.Round(in / l.scaleFactor)
}

func skipNodes(start int, nodes []nodeT) int {
	for j, node := range nodes[start:] {
		// skip nodes after a line break
		if node.isBox() || node.isForcedBreak() {
			start += j

			break
		}
	}

	return start
}

// render the nodes with the provided line breaks.
//
// Rendering is for now essentially for a plain terminal output, with fixed-width fonts.
func (l *LineBreaker) render(breakList *breakPoint) []string {
	var (
		lines  []lineT
		result []string
	)

	lineStart := 0
	for brk := breakList.next; brk != nil; brk = brk.next {
		lineStart = skipNodes(lineStart, l.nodes)

		lines = append(lines, lineT{
			ratio:    brk.ratio,
			nodes:    l.nodes[lineStart : brk.position+1],
			position: brk.position,
		})

		lineStart = brk.position
	}
	attributesState := l.attrList.Iterator() // TODO: First()

	for _, line := range lines {
		lineResult := new(strings.Builder)
		runesWriter := runesio.AsRunesWriter(lineResult)
		attributesState.StartOfLine(runesWriter)
		for index, node := range line.nodes {
			switch {
			case node.isBox():
				// render a box node
				node.Render(runesWriter)
				if node.HasRenderer() { // TODO: add state handling to box node
					_ = attributesState.Next()
				}

			case node.isGlue():
				// render a glue node
				pad := node.width
				if line.ratio < 0 {
					pad += line.ratio * node.shrink
				} else {
					pad += line.ratio * node.stretch
				}

				spaces := repeatRunes(space, int(l.downScale(pad)))
				_, _ = runesWriter.WriteRunes(spaces)

			case node.isPenalty():
				if l.renderHyphens && node.penalty == l.hyphenPenalty && index == len(line.nodes)-1 {
					// render a soft hyphen node
					_, _ = runesWriter.WriteRunes(hyphen)
				}
			}
		}
		attributesState.EndOfLine(runesWriter)

		result = append(result, lineResult.String())
	}

	return result
}

func repeatRunes(in []rune, times int) []rune {
	if len(in) == 0 || times == 0 {
		return []rune{}
	}

	out := make([]rune, 0, len(in)*times)
	for j := 0; j < times; j++ {
		out = append(out, in...)
	}

	return out
}

// boxNodes models box nodes with possible word breaks (suited for left-aligned output).
//
// Raw tokens are split into:
// * Renderers with the appropriate start/end ANSI control sequence
// * word parts separated by punctuation marks and other separators (not hyphens)
// * word parts at legit hyphenation breakpoints
func (l *LineBreaker) boxNodes(token []rune) []nodeT {
	nodes := make([]nodeT, 0, 10)

	for _, stripped := range ansi.StripToken(token) { // there may be several start/stop escape sequences: break them down
		tokenState := newTokenState(stripped, l.attrList)

		// this text has been stripped from start/stop escape sequences. The attribute renderer will remember the start/stop sequences.
		// We don't necessarily need to create as many renderers, but we must keep track of the state
		for _, strippedFromPunct := range l.punctuator(stripped.Text) { // split punctuation marks as well as separators such as "/", "|", "_"...
			if punctuator.IsPunctuation(strippedFromPunct) {
				tokenState.Start(strippedFromPunct)

				// A punctuation mark, or similar separator (e.g. "/", "|", "&"...).
				//
				// NOTE(fredbi): nice to have - we might want to distinguish different rules depending on
				// the punctuation mark. E.g. ";", "&", "." should probably deserve a special processing.
				// For the moment, this package essentially supports rendering with fixed-width fonts on a terminal, so this is not really needed.
				nodes = append(nodes,
					newPenalty(noWidth, infinity, unflaggedPenalty),
					newGlue(noWidth, l.glueStretch, noShrink), // no space before the punctuation mark. Some typographic rules disagree with this, e.g. for ";"
					newBox(l.scale(l.measurer(strippedFromPunct)), strippedFromPunct, tokenState.Current()),
					newPenalty(noWidth, l.punctuationPenalty, flaggedPenalty), // this penalty won't be mixed with hyphens
					newGlue(noWidth, -l.glueStretch, noShrink),
				)

				continue
			}

			if !l.wordBreak || len(strippedFromPunct) <= l.minHyphenate {
				// Either word breaking is forbidden or this token is too short for a legitimate hyphenation
				tokenState.Start(strippedFromPunct)
				nodes = append(nodes, newBox(l.scale(l.measurer(strippedFromPunct)), strippedFromPunct, tokenState.Current()))

				continue
			}

			for _, word := range hyphenator.SplitWord(strippedFromPunct) { // split on explicit hyphens
				// An explicit hyphen: this will be rendered as a regular token, but provides a legit line break point.
				if hyphenator.IsHyphen(word) {
					tokenState.Start(word)
					nodes = append(nodes,
						newPenalty(noWidth, infinity, unflaggedPenalty),
						newGlue(noWidth, l.glueStretch, noShrink),
						newBox(l.scale(l.measurer(word)), word, tokenState.Current()),
						newPenalty(noWidth, l.hardHyphenPenalty, flaggedPenalty), // this penalty won't be mixed with soft hyphens
						newGlue(noWidth, -l.glueStretch, noShrink),
					)

					continue
				}

				// Soft hyphens: the hyphenator returns word parts, broken at legit hyphenation breakpoints
				hyphenated := l.hyphenator(word)

				// word break points are associated with a penalty
				for _, part := range hyphenated[:len(hyphenated)-1] {
					tokenState.Start(part)
					nodes = append(nodes,
						newBox(l.scale(l.measurer(part)), part, tokenState.Current()),
					)
					nodes = append(nodes, l.pushHyphen()...)
				}

				lastPart := hyphenated[len(hyphenated)-1]
				tokenState.Start(lastPart)
				nodes = append(nodes,
					newBox(l.scale(l.measurer(lastPart)), lastPart, tokenState.Current()),
				)
			}
		}
		// add stop to the last renderer
		tokenState.Stop()
	}

	return nodes
}

func (l *LineBreaker) pushHyphen() []nodeT {
	if l.renderHyphens {
		// when rendering hyphens, the penalty incurs some consumed width
		return []nodeT{
			// justified:
			// newPenalty(l.hyphenWidth, l.hyphenPenalty, flaggedPenalty),
			// ragged right:
			newPenalty(noWidth, infinity, unflaggedPenalty),
			newGlue(noWidth, l.glueStretch, noShrink),
			newPenalty(l.hyphenWidth, l.hyphenPenalty, flaggedPenalty),
			newGlue(noWidth, -l.glueStretch, noShrink),
		}
	}

	return []nodeT{
		// when hyphens are not rendered (words are just broken), there is no width associated to the penalty
		newPenalty(noWidth, l.hyphenPenalty, flaggedPenalty),
	}
}

// leftAlignedNodes prepares nodes for left-aligned rendering (ragged right).
func (l *LineBreaker) leftAlignedNodes(tokens []string) []nodeT {
	if len(tokens) == 0 {
		return nil
	}

	nodes := make([]nodeT, 0, 4*(len(tokens)-1)+3)

	// transform tokens into a list of nodes of type (box|glue|penalty)
	for _, word := range tokens[:len(tokens)-1] {
		nodes = append(nodes, l.boxNodes([]rune(word))...) // a word token, possibly broken in parts
		// from K&P: justified:
		// nodes = append(nodes, newGlue(l.spaceWidth, l.glueStretch, l.glueShrink))
		// from K&P: ragged right:
		nodes = append(nodes, newGlue(noWidth, l.glueStretch, noShrink))
		nodes = append(nodes, newPenalty(noWidth, 0, unflaggedPenalty))
		nodes = append(nodes, newGlue(l.spaceWidth, -l.glueStretch, noShrink))
	}

	// last token: complete the list of nodes with a final infinite glue and penalty.
	nodes = append(nodes, l.boxNodes([]rune(tokens[len(tokens)-1]))...)
	nodes = append(nodes, newGlue(noWidth, infinity, noShrink))
	nodes = append(nodes, newPenalty(noWidth, -infinity, flaggedPenalty))

	return nodes
}

// TODO: center(), justify()?
