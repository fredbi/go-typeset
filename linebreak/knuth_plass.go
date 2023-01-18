package linebreak

import (
	"container/list"
	"math"
)

const (
	infinity   = 10000.00
	maxDemerit = math.MaxFloat64
)

func (l *LineBreaker) adjustmentRatio(fromBreakPoint *breakPoint, toNode nodeT, sum *sums, idealWidth float64) float64 {
	actualWidth := sum.width - fromBreakPoint.totals.width

	if toNode.isPenalty() {
		// for penalties with a width, e.g. extra hyphen
		actualWidth += toNode.width
	}

	switch {
	case actualWidth < idealWidth:
		// need to stretch
		stretch := sum.stretch - fromBreakPoint.totals.stretch

		if stretch > 0 {
			return (idealWidth - actualWidth) / stretch
		}

		return infinity

	case actualWidth > idealWidth:
		// need to shrink
		shrink := sum.shrink - fromBreakPoint.totals.shrink

		if shrink > 0 {
			return (idealWidth - actualWidth) / shrink
		}

		return infinity

	default:
		// perfect match
		return 0.00
	}
}

// demeritsForRatio attributes a score to the adjustment ratio, taking penalties into account.
func (l *LineBreaker) demeritsForRatio(node nodeT, ratio float64) float64 {
	badness := l.demerits.line + l.badness*math.Pow(math.Abs(ratio), 3)
	penalty := node.penalty

	switch {
	case node.isPenalty() && penalty > 0:
		return math.Pow(badness+penalty, 2)
	case node.isPenalty() && penalty != -infinity:
		return math.Pow(badness, 2) - math.Pow(penalty, 2)
	default:
		return math.Pow(badness, 2)
	}
}

// demeritsAndClass attributes a demerits score and fitness class from a break point to a node.
func (l *LineBreaker) demeritsAndClass(fromBreakPoint *breakPoint, toNode nodeT, ratio float64) (float64, fitnessClass) {
	demerits := l.demeritsForRatio(toNode, ratio)

	previous := l.nodes[fromBreakPoint.position]
	if toNode.isPenalty() && toNode.flagged && previous.isPenalty() && previous.flagged {
		// penalize consecutive flagged penalty nodes
		demerits += l.demerits.flagged
	}

	currentClass := newFitnessClass(ratio)

	// add demerits due to fitness class whenever the fitness of 2 adjacent lines differ too much
	if currentClass.isAwayFrom(fromBreakPoint.fitness) {
		demerits += l.demerits.fitness
	}

	demerits += fromBreakPoint.totalDemerits

	return demerits, currentClass
}

// breakPoints yields an ordered linked-list breakpoints
func (l *LineBreaker) breakPoints() *breakPoint {
	// reset state
	l.sum = new(sums)
	l.activeNodes = list.New()
	l.activeNodes.PushBack(newBreakPoint(0, 0, 0, 0, fitnessClassOne, sums{}, nil)) // first empty node starting a paragraph
	startNode := l.findStartNode()                                                  // Baskerville version - should be 1 in normal cases

	for i, node := range l.nodes[startNode:] {
		index := startNode + i

		switch {
		case node.isBox():
			l.sum.width += node.width // accumulate the total width of word
		case node.isGlue():
			if index > 0 && l.nodes[index-1].nodeType == nodeTypeBox {
				l.mainLoop(index) // explore a glue following a word
			}

			l.sum.Add(node.sums)

		case node.isPenalty() && node.penalty != infinity:
			l.mainLoop(index) // explore a penalty
		}
	}

	if l.activeNodes.Len() > 0 {
		nodeWithMinDemerits := l.findBestBreak()

		return reverseBreakPoints(nodeWithMinDemerits)
	}

	if l.looseness == 0 {
		return nil
	}

	// choose appropriate node
	// TODO: implem looseness != 0 ("choose the appropriate active node")

	return nil
}

// findStartNode skips starting glues (i.e. indentations) or penalties.
// TODO: remove?? (Baskerville version only)
func (l *LineBreaker) findStartNode() (start int) {
	for _, node := range l.nodes {
		switch node.nodeType {
		case nodeTypeBox:
			return start

		case nodeTypePenalty:
			if node.penalty == -infinity {
				return start
			}
		default:
			start++
		}
	}

	return start
}

func (l *LineBreaker) findBestBreak() *breakPoint {
	nodeWithMinDemerits := &breakPoint{
		totalDemerits: maxDemerit,
	}

	for element := l.activeNodes.Front(); element != nil; element = element.Next() {
		node := element.Value.(*breakPoint)

		if node.totalDemerits < nodeWithMinDemerits.totalDemerits {
			nodeWithMinDemerits = node
		}
	}

	return nodeWithMinDemerits
}

func (l *LineBreaker) sumFromNode(index int) sums {
	sum := *l.sum

	for i, node := range l.nodes[index:] {
		if node.isGlue() {
			sum.Add(node.sums)

			continue
		}

		if node.isBox() || (node.isForcedBreak() && i > 0) {
			break
		}
	}

	return sum
}

// exploreForNode is referred to as "the main loop" in Knuth & Plass.
func (l *LineBreaker) mainLoop(index int) {
	node := l.nodes[index]
	activeElement := l.activeNodes.Front()

	var (
		currentLine int // will range over lines starting from 1
		candidates  map[fitnessClass]*breakPoint
	)

	lowestDemerits := maxDemerit
	for activeElement != nil {
		candidates = defaultCandidates() // set candidates with infinite demerits

		// break points up to the current line
		for activeElement != nil {
			active := activeElement.Value.(*breakPoint)
			next := activeElement.Next()
			currentLine = active.line + 1
			ratio := l.adjustmentRatio(active, node, l.sum, l.idealWidth(currentLine))

			if ratio < -1 || node.isForcedBreak() {
				// deactivate an undesirable break or a forced line break
				l.activeNodes.Remove(activeElement)
			}

			if ratio >= -1 && ratio <= l.tolerance {
				// update candidate
				demerits, currentClass := l.demeritsAndClass(active, node, ratio)
				lowestDemerits = minf(lowestDemerits, demerits)

				if demerits < candidates[currentClass].totalDemerits {
					candidates[currentClass] = active
				}
			}

			// ratio > l.tolerance is not considered feasible

			activeElement = next

			if activeElement != nil && active.line >= currentLine {
				// stop iterating to add new candidates
				break
			}
		}

		if lowestDemerits < maxDemerit {
			l.insertNewActiveBreak(activeElement, index, lowestDemerits, candidates)
		}
	}
}

func (l *LineBreaker) insertNewActiveBreak(activeElement *list.Element, index int, lowestDemerits float64, candidates map[fitnessClass]*breakPoint) {
	sum := l.sumFromNode(index)

	for class, candidate := range candidates {
		if candidate.totalDemerits >= maxDemerit || candidate.totalDemerits > lowestDemerits+l.demerits.fitness {
			// skip default candidate
			continue
		}

		newBreak := newBreakPoint(
			index,                                    // break at node index
			candidate.totalDemerits, candidate.ratio, // ratings for this break point
			candidate.line+1,
			class,
			sum,       // totals after this node
			candidate, // link to the previous candidate breakpoint
		)

		if activeElement != nil {
			_ = l.activeNodes.InsertBefore(newBreak, activeElement)

			return
		}

		_ = l.activeNodes.PushBack(newBreak)
	}
}

// reverseBreakPoints walks backwards the list of breakpoints,
// and prepares the ordered (forward) walking.
func reverseBreakPoints(brk *breakPoint) *breakPoint {
	var first *breakPoint
	brk.next = nil

	for brk != nil {
		if brk.previous != nil {
			brk.previous.next = brk
			first = brk.previous
		}

		brk = brk.previous
	}

	return first
}

// getIdealWidth retrieves the constraint on the line length.
//
// NOTE: currentLine starts at 1.
func (l *LineBreaker) idealWidth(currentLine int) float64 {
	if currentLine < len(l.lineWidths)+1 {
		return l.lineWidths[currentLine-1]
	}

	return l.lineWidths[len(l.lineWidths)-1]
}
