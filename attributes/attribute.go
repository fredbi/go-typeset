package attributes

import (
	"github.com/fredbi/go-typeset/terminal/runes/runesio"
)

type (
	// Attribute implements a basic renderer.
	//
	// The start and stop mode are rendered simply by writing their runes.
	// This typically corresponds to a terminal use-case with ANSI escape sequences.
	Attribute struct {
		Text          []rune
		StartSequence []rune
		StopSequence  []rune
		NestingLevel  int
	}
)

// New attribute that will output to a io.Writer.
func New(text, start, stop []rune) *Attribute {
	return &Attribute{
		Text:          text,
		StartSequence: start,
		StopSequence:  stop,
	}
}

func (a Attribute) Start(w runesio.Writer) {
	if a.HasStart() {
		_, _ = w.WriteRunes(a.StartSequence)
	}
}

func (a Attribute) Stop(w runesio.Writer) {
	if a.HasStop() {
		_, _ = w.WriteRunes(a.StopSequence)
	}
}

func (a Attribute) Runes() []rune {
	return a.Text
}

func (a Attribute) Level() int {
	return a.NestingLevel
}

func (a *Attribute) SetLevel(level int) {
	a.NestingLevel = level
}

func (a *Attribute) SetStop(stop []rune) {
	a.StopSequence = stop
}

func (a Attribute) HasStart() bool {
	return len(a.StartSequence) > 0
}

func (a Attribute) HasStop() bool {
	return len(a.StopSequence) > 0
}

// Render an attribute.
func (a Attribute) Render(w runesio.Writer) {
	a.Start(w)
	if len(a.Text) > 0 {
		_, _ = w.WriteRunes(a.Text)
	}
	a.Stop(w)
}
