package linebreak

/*
// attention: currentLine starts with 1
// ET LA
func (l *LineBreaker) costRatio2(sum *sums, index int, previous *sums, currentLine int) float64 {
	actualWidth := sum.width - previous.width
	idealWidth := getIdealWidth(currentLine, l.lineLengths)

	if l.nodes[index].nodeType == nodeTypePenalty {
		// for penalties with a width, e.g. extra hyphen
		actualWidth += l.nodes[index].width
	}

	switch {
	case actualWidth < idealWidth:
		// need to stretch
		stretch := sum.stretch - previous.stretch

		if stretch > 0 {
			return float64(idealWidth-actualWidth) / float64(stretch)
		}

		return infinity

	case actualWidth > idealWidth:
		// need to shrink
		shrink := sum.shrink - previous.shrink

		if shrink > 0 {
			return float64(idealWidth-actualWidth) / float64(shrink)
		}

		return infinity

	default:
		// perfect match
		return 0.00
	}
}
*/
/*
func (l *LineBreaker) fallback(nodes []nodeT, lineLengths []int) []breakPoint {
	l.lineLengths = lineLengths
	l.nodes = nodes
	l.sum = new(sums)
	l.activeNodes = list.New()
	l.activeNodes.PushBack(newBreakPoint(0, 0, 0, 0, 0, sums{}, nil)) // first empty node starting a paragraph

	startNode := l.findStartNode()
	currentLine := 1
	previousIndex := 0
	previousSum := *l.sum
	consecutiveFlagged := 0

MAINLOOP:
	for index := startNode; index < len(l.nodes); index++ {
		current := l.nodes[index]

		switch current.nodeType {
		case nodeTypeBox:
			l.sum.width += current.width
			idealWidth := getIdealWidth(currentLine, l.lineLengths)

			if l.sum.width-previousSum.width > idealWidth {
				ratio := l.costRatio2(l.sum, index, &previousSum, currentLine)

				if ratio < -1 {
					highIndex := index
					highSums := *l.sum
					boxes := 0

				LOOP:
					for ; index > previousIndex; index-- {
						current = l.nodes[index]

						switch current.nodeType {
						case nodeTypeBox:
							l.sum.width -= current.width
							boxes++
						case nodeTypeGlue:
							l.sum.Minus(current.sums)
						case nodeTypePenalty:
							break LOOP
						}
					}

					if index == previousIndex {
						if boxes == highIndex-previousIndex {
							return []breakPoint{} // fail?
						}

						index = highIndex
						l.sum = &highSums

						for ; index > previousIndex; index-- {
							if l.nodes[index].nodeType == nodeTypeBox {
								break
							}

							l.sum.width += l.nodes[index].width
						}

						ratio := l.costRatio2(l.sum, index, &previousSum, currentLine)
						if ratio > l.tolerance {
							lowIndex := index
							lowSums := *l.sum
							lowRatio := ratio
							// lowWidth := l.sum.width - previousSum.width
							index := highIndex
							*l.sum = highSums

						LOOP2:
							for ; index > lowIndex; index-- {
								current = l.nodes[index]

								switch current.nodeType {
								case nodeTypeBox:
									l.sum.width -= current.width
								case nodeTypeGlue:
									l.sum.Minus(current.sums)
								case nodeTypePenalty:
									if current.penalty < infinity && (current.flagged == 0 || consecutiveFlagged < 2) {
										ratio = l.costRatio2(l.sum, index, &previousSum, currentLine)
										if ratio > -1 && ratio <= lowRatio {
											break LOOP2
										}
									}
								}
							}

							if index == lowIndex {
								*l.sum = lowSums
								ratio = lowRatio
							}
						}
					} else {

						if index == len(l.nodes)-1 || l.nodes[index+1].nodeType != nodeTypeGlue {
							continue MAINLOOP
						}

						index++ // skip
					}

					previousIndex =  index
					previousSums = *l.sums

					l.activeNodes.PushBack(newBreakPoint(
				index,                               // break at node index
				0, ratio, // ratings for this break point
				candidate.active.line+1, fitnessClass,
				sum,              // totals from the node
				candidate.active, // link to the previous candidate breakpoint
			)
						// TODO
					))
				}
			}
		case nodeTypeGlue:
			l.sum.Add(current.sums)
		case nodeTypePenalty:
		}
	}
}
*/
