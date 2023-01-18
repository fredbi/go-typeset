package linebreak

import (
	"github.com/fredbi/go-typeset/attributes"
	"github.com/fredbi/go-typeset/terminal/ansi"
	"github.com/fredbi/go-typeset/terminal/runes/runesio"
)

type (
	nodeT struct {
		nodeType  nodeType
		penalty   float64
		flagged   bool
		value     []rune
		attribute attributes.Renderer // attributes such as color, italic, bold ...

		sums
	}

	sums struct {
		width   float64 // width
		stretch float64 // for glues, stretchability parameter
		shrink  float64 // for glues, shrinkability parameter
	}

	breakPoint struct {
		position      int          // index of this breakpoint (0 = start of paragraph)
		line          int          // the line ending at this breakpoint
		fitness       fitnessClass // fitness class of the line ending at this breakpoint
		totalDemerits float64      // minimum total demerits up to this breakpoint
		totals        sums         // total width, stretch and shrink used to calculate adjustment ratios
		ratio         float64      // informative: the adjustment ratio at this breakpoint
		previous      *breakPoint  // pointer to the best node for the preceding breakpoint
		next          *breakPoint
	}

	nodeType uint8

	demeritsT struct {
		line    float64
		flagged float64
		fitness float64
	}

	lineT struct {
		ratio    float64
		nodes    []nodeT
		position int
	}

	fitnessClass int

	err string

	tokenState struct {
		list            *attributes.State
		currentRenderer attributes.Renderer
		isStarted       bool
		isStopped       bool
		stripped        ansi.StrippedToken
	}
)

const (
	// ErrCannotBeSet indicates that the display constraints cannot be met with the given tolerance parameter.
	ErrCannotBeSet err = "paragraph cannot be set with the given tolerance"
)

const (
	// node types

	nodeTypePenalty nodeType = iota
	nodeTypeGlue
	nodeTypeBox
)

const (
	// fitness classes

	fitnessClassZero  fitnessClass = iota // "tight lines"
	fitnessClassOne                       // "normal lines"
	fitnessClassTwo                       // "loose lines"
	fitnessClassThree                     // "very loose lines"
)

func newBreakPoint(position int, demerits float64, ratio float64, line int, fitness fitnessClass, totals sums, previous *breakPoint) *breakPoint {
	return &breakPoint{
		position:      position,
		totalDemerits: demerits,
		ratio:         ratio,
		line:          line,
		fitness:       fitness,
		totals:        totals,
		previous:      previous,
	}
}

//nolint:unparam
func newGlue(width, stretch, shrink float64) nodeT {
	return nodeT{
		nodeType: nodeTypeGlue,
		sums: sums{
			width:   width,
			stretch: stretch,
			shrink:  shrink,
		},
	}
}

func newBox(width float64, value []rune, attribute attributes.Renderer) nodeT {
	return nodeT{
		nodeType: nodeTypeBox,
		value:    value,
		sums: sums{
			width: width,
		},
		attribute: attribute,
	}
}

func newPenalty(width float64, penalty float64, flagged bool) nodeT {
	return nodeT{
		nodeType: nodeTypePenalty,
		sums: sums{
			width: width,
		},
		penalty: penalty,
		flagged: flagged,
	}
}

// newFitnessClass establish a coarse fitness classification according to the cost ratio.
func newFitnessClass(ratio float64) fitnessClass {
	switch {
	case ratio < -0.5:
		return fitnessClassZero
	case ratio <= 0.5:
		return fitnessClassOne
	case ratio <= 1:
		return fitnessClassTwo
	default:
		return fitnessClassThree
	}
}

func defaultCandidates() map[fitnessClass]*breakPoint {
	return map[fitnessClass]*breakPoint{
		fitnessClassZero:  newBreakPoint(0, maxDemerit, 0, 0, fitnessClassZero, sums{}, nil),
		fitnessClassOne:   newBreakPoint(0, maxDemerit, 0, 0, fitnessClassOne, sums{}, nil),
		fitnessClassTwo:   newBreakPoint(0, maxDemerit, 0, 0, fitnessClassTwo, sums{}, nil),
		fitnessClassThree: newBreakPoint(0, maxDemerit, 0, 0, fitnessClassThree, sums{}, nil),
	}
}

func (e err) Error() string {
	return string(e)
}

func (s *sums) Add(t sums) {
	s.width += t.width
	s.stretch += t.stretch
	s.shrink += t.shrink
}

func (t nodeType) String() string {
	switch t {
	case nodeTypePenalty:
		return "penalty"
	case nodeTypeGlue:
		return "glue"
	case nodeTypeBox:
		return "box"
	default:
		return ""
	}
}

func (n nodeT) isGlue() bool {
	return n.nodeType == nodeTypeGlue
}

func (n nodeT) isBox() bool {
	return n.nodeType == nodeTypeBox
}

func (n nodeT) isPenalty() bool {
	return n.nodeType == nodeTypePenalty
}

func (n nodeT) isForcedBreak() bool {
	return n.nodeType == nodeTypePenalty && n.penalty == -infinity
}

func (n nodeT) Render(w runesio.Writer) {
	if !n.isBox() {
		return
	}

	if n.attribute != nil {
		n.attribute.Render(w)

		return
	}

	_, _ = w.WriteRunes(n.value)
}

func (n nodeT) HasRenderer() bool {
	return n.attribute != nil
}

func (f fitnessClass) isAwayFrom(g fitnessClass) bool {
	return abs(int(f)-int(g)) > 1
}

func (s *tokenState) Start(text []rune) {
	if !s.isStarted {
		s.isStarted = true

		s.currentRenderer = attributes.New(text, s.stripped.StartSequence, nil)
	} else {
		s.currentRenderer = attributes.New(text, nil, nil)
	}

	s.list.Push(s.currentRenderer)
}

func (s *tokenState) Stop() {
	last := s.currentRenderer
	if last == nil || len(s.stripped.StopSequence) == 0 {
		return
	}

	last.SetStop(s.stripped.StopSequence)
}

func (s *tokenState) Current() attributes.Renderer {
	return s.currentRenderer
}

func newTokenState(stripped ansi.StrippedToken, list *attributes.State) *tokenState {
	return &tokenState{
		list:     list,
		stripped: stripped,
	}
}

/*
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
*/

func minf(a, b float64) float64 {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}

	return a
}
