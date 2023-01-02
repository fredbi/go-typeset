package linebreak

import (
	"container/list"
	"math"
	"strings"
)

type (
	// LineBreaker breaks a paragraph into lines under line width constraints.
	//
	// It implements the Knuth-Plass algorithm.
	LineBreaker struct {
		nodes       []nodeT
		lineWidths  []float64
		sum         *sums
		activeNodes *list.List
		spaceWidth  float64
		hyphenWidth float64
		// spaceStretch float64
		//spaceShrink  float64

		*options
	}
)

const (
	infinity   = 10000.00
	maxDemerit = math.MaxFloat64

	hyphen           = "-"
	space            = " "
	flaggedPenalty   = true
	unflaggedPenalty = false
	noWidth          = 0.0
	noShrink         = 0.0
)

func New(opts ...Option) *LineBreaker {
	l := &LineBreaker{
		options: defaultOptions(opts),
	}

	l.spaceWidth = l.scale(l.measurer(space))
	if l.wordBreak {
		l.hyphenWidth = l.scale(l.measurer(hyphen))
	}

	// NOTE: this is for center & justify (not implemented for now)
	// l.spaceStretch = min(1, int(float64(l.spaceWidth*l.space.width)/float64(l.space.stretch)))
	// l.spaceShrink = min(1, int(float64(l.spaceWidth*l.space.width)/float64(l.space.shrink)))

	return l
}

// LeftAlignUniform left-align a series of tokens that compose a paragraph,
// rendering multiple lines of uniform length maxLength.
func (l *LineBreaker) LeftAlignUniform(tokens []string, maxLength float64) ([]string, error) {
	l.nodes = l.leftAlignedNodes(tokens)
	l.lineWidths = l.buildUniformLengths(maxLength)

	breakList := l.breakPoints()
	if breakList == nil {
		return nil, ErrCannotBeSet
	}

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
	// return int(math.Round(float64(in) * r / l.scaleFactor))
	return math.Round(float64(in) / l.scaleFactor)
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

// render the nodes with the provided line breaks .
//
// TODO: state handling for ANSI escape sequences (box attributes)
// TODO: iterate breaks from list directly
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

	for _, line := range lines {
		lineResult := new(strings.Builder)

		for index, node := range line.nodes {
			switch {
			case node.isBox():
				// render a box node
				lineResult.WriteString(node.value)

			case node.isGlue():
				// render a glue node
				pad := node.width
				if line.ratio < 0 {
					pad += line.ratio * node.shrink
				} else {
					pad += line.ratio * node.stretch
				}

				lineResult.WriteString(strings.Repeat(space, int(l.downScale(pad))))

			case node.isPenalty():
				if l.renderHyphens && node.penalty == l.hyphenPenalty && index == len(line.nodes)-1 {
					// render a hyphen penalty node
					lineResult.WriteString(hyphen)
				}
			}
		}

		result = append(result, lineResult.String())
	}

	return result
}

// boxNodes produces box nodes with possible word breaks (suited for left-aligned output).
// TODO(fredbi): process punctuation marks, word breaking on natural separators (e.g. /|-_)
func (l *LineBreaker) boxNodes(word string) []nodeT {
	if !l.wordBreak || len(word) <= l.minHyphenate {
		return []nodeT{newBox(l.scale(l.measurer(word)), word)}
	}

	hyphenated := l.hyphenator(word)
	if len(hyphenated) == 0 { // this rule is based on the # runes, not the width
		return []nodeT{newBox(l.scale(l.measurer(word)), word)}
	}
	// TODO: explicit hyphens or dashes

	boxNodes := make([]nodeT, 0, len(hyphenated))

	// word break points are associated with a penalty
	for _, part := range hyphenated[:len(hyphenated)-1] {
		boxNodes = append(boxNodes,
			newBox(l.scale(l.measurer(part)), part),
		)

		if l.renderHyphens {
			// when rendering hyphens, the penalty incurs some consumed width
			boxNodes = append(boxNodes,
				// justified:
				//newPenalty(l.hyphenWidth, l.hyphenPenalty, flaggedPenalty),
				// ragged right:
				newPenalty(noWidth, infinity, unflaggedPenalty),
				newGlue(noWidth, l.glueStretch, noShrink),
				newPenalty(l.hyphenWidth, l.hyphenPenalty, flaggedPenalty),
				newGlue(noWidth, -l.glueStretch, noShrink),
			)
		} else {
			// when hyphens are not rendered (words are just broken), there is no width associated to the penalty
			boxNodes = append(boxNodes,
				newPenalty(noWidth, l.hyphenPenalty, flaggedPenalty),
			)
		}
	}

	lastPart := hyphenated[len(hyphenated)-1]
	boxNodes = append(boxNodes,
		newBox(l.scale(l.measurer(lastPart)), lastPart),
	)

	return boxNodes
}

// leftAlignedNodes prepares nodes for left-alignment.
// TODO: punctuation marks
func (l *LineBreaker) leftAlignedNodes(tokens []string) []nodeT {
	nodes := make([]nodeT, 0, 4*(len(tokens)-1)+3)

	// transform tokens into a list of nodes of type (box|glue|penalty)
	for _, word := range tokens[:len(tokens)-1] {
		nodes = append(nodes, l.boxNodes(word)...) // a word token, possibly broken in parts
		// from K&P: justified:
		// nodes = append(nodes, newGlue(l.spaceWidth, l.glueStretch, l.glueShrink))
		// from K&P: ragged right:
		nodes = append(nodes, newGlue(noWidth, l.glueStretch, noShrink))
		nodes = append(nodes, newPenalty(noWidth, 0, unflaggedPenalty))
		nodes = append(nodes, newGlue(l.spaceWidth, -l.glueStretch, noShrink))
	}

	// last token: complete the list of nodes with a final infinite glue and penalty.
	nodes = append(nodes, l.boxNodes(tokens[len(tokens)-1])...)
	nodes = append(nodes, newGlue(noWidth, infinity, noShrink))
	nodes = append(nodes, newPenalty(noWidth, -infinity, flaggedPenalty))

	return nodes
}

// TODO: center(), justify()?
