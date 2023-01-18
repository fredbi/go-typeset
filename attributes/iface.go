package attributes

import (
	"github.com/fredbi/go-typeset/terminal/runes/runesio"
)

type (
	// Renderers is a collection of Renderer(s)
	Renderers []Renderer

	// Renderer knows how to
	Renderer interface {
		// Start renders the start formatting sequence, if any.
		Start(runesio.Writer)
		// Start renders the stop formatting sequence, if any.
		Stop(runesio.Writer)
		// Runes renders the runes of the token, without formatting sequences.
		Runes() []rune
		// Level indicates the level of nested start/stop sequences.
		Level() int
		// SetLevel sets the level of nested start/stop sequences.
		SetLevel(int)
		// HasStart is true when the renderer contains a start sequence.
		HasStart() bool
		// HasStop is true when the renderer contains a stop sequence.
		HasStop() bool

		// Render knows how to render the token with the formatting sequences.
		Render(runesio.Writer)

		SetStop([]rune)

		// Iterator walks a list of renderers from the current element.
		// Iterator() Iterator
	}

	// Iterator walks a chained list of Renderers.
	Iterator interface {
		// Next indicates if there is a next item to be consumed.
		Next() bool
		// Item returns the currently iterated Renderer.
		Item() Renderer
	}
)

func (c Renderers) Start(w runesio.Writer) {
	for _, r := range c {
		if r == nil {
			continue
		}

		r.Start(w)
	}
}

func (c Renderers) Stop(w runesio.Writer) {
	for _, r := range c {
		if r == nil {
			continue
		}

		r.Stop(w)
	}
}

func (c Renderers) Render(w runesio.Writer) {
	for _, r := range c {
		if r == nil {
			continue
		}

		r.Render(w)
	}
}
