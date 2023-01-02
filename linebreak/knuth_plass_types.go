package linebreak

type (
	nodeT struct {
		nodeType   nodeType
		penalty    float64
		flagged    bool
		value      string
		attributes []string // attributes such as color, italic, bold ...

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

func newBox(width float64, value string) nodeT {
	return nodeT{
		nodeType: nodeTypeBox,
		value:    value,
		sums: sums{
			width: width,
		},
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

func (f fitnessClass) isAwayFrom(g fitnessClass) bool {
	return abs(int(f)-int(g)) > 1
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

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
