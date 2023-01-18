package attributes

import (
	"container/list"

	"github.com/fredbi/go-typeset/terminal/runes/runesio"
)

type (
	// State maintains the state of attribute renderers.
	//
	// The state is constructed from a flow of Renderers using Push().
	State struct {
		list *list.List
	}

	// StateIterator iterates over the list of renderers.
	//
	// A StateIterator knows how to close the current rendering context when reaching the end of line,
	// and to resume the previous rendering context when starting a new line.
	StateIterator struct {
		start   bool
		current *list.Element
	}
)

func NewState() *State {
	return &State{
		list: list.New(),
	}
}

func (s *State) Push(attr Renderer) {
	if attr == nil {
		return
	}
	current := s.list.Back()

	// determine the nesting level of attributes for this element
	var currentLevel int
	if current != nil {
		attribute := current.Value.(Renderer)
		currentLevel = attribute.Level()
		hadStop := attribute.HasStop()
		if hadStop {
			currentLevel--
		}
	}

	if attr.HasStart() {
		currentLevel++
	}

	attr.SetLevel(currentLevel)

	_ = s.list.PushBack(attr)
}

// Iterator yields an iterator to walk over the renderers in this State, from first to last.
func (s *State) Iterator() *StateIterator {
	return &StateIterator{
		start:   true,
		current: s.list.Front(),
	}
}

func (s *State) First() *StateIterator {
	return &StateIterator{
		start:   false,
		current: s.list.Front(),
	}
}

// Next yields true if there are more items to consume.
func (i *StateIterator) Next() bool {
	if i.start {
		i.start = false

		return i.current != nil
	}

	if i.current == nil {
		return false
	}

	i.current = i.current.Next()

	return i.current != nil
}

// Item returns the currently iterated Renderer.
func (i *StateIterator) Item() Renderer {
	if i.current == nil {
		return nil
	}

	return i.current.Value.(Renderer)
}

func (i *StateIterator) backtrack(level int) *list.Element {
	if i.current == nil {
		return nil
	}

	current := i.current
	/*
		attribute := current.Value.(Renderer)
		currentLevel := attribute.Level()
	*/

	// backtrack to the first renderer for this level
	for current != nil {
		lookback := current.Prev()
		if lookback == nil {
			return current
		}

		attribute := lookback.Value.(Renderer)
		if attribute.Level() < level {
			return current
		}

		current = lookback
	}

	return nil
}

// StartOfLine restores the state of attributes to start a new line.
func (i *StateIterator) StartOfLine(w runesio.Writer) {
	if i.current == nil {
		return
	}

	// if the end of line was preceded by a stop, go down 1 level
	attribute := i.current.Value.(Renderer)
	currentLevel := attribute.Level()
	if attribute.HasStop() {
		currentLevel--
	}
	// backtrack to the first renderer for this level
	current := i.backtrack(currentLevel)
	if current == nil || currentLevel == 0 {
		return
	}

	// starts all rendering sequences for this level
	for current != nil {
		attribute := current.Value.(Renderer)

		attribute.Start(w)

		if current == i.current {
			break
		}

		current = current.Next()
	}
}

// EndOfLine closes the state of attributes to end the line with an empty state.
func (i *StateIterator) EndOfLine(w runesio.Writer) {
	if i.current == nil {
		return
	}

	// if the end of line was preceded by a stop, go down 1 level
	attribute := i.current.Value.(Renderer)
	currentLevel := attribute.Level()
	if attribute.HasStop() {
		currentLevel--
	}

	// backtrack to the first renderer for this level
	current := i.backtrack(currentLevel)
	if current == nil || currentLevel == 0 {
		return
	}

	// stop all sequences up to the last renderer at this level
	for current != nil {
		attribute := current.Value.(Renderer)
		level := attribute.Level()

		if level < currentLevel {
			break
		}

		attribute.Stop(w)

		current = current.Next()
	}
}
