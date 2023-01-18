package linebreak

import (
	"github.com/fredbi/go-typeset/terminal/runes"
	wordbreaker "github.com/fredbi/go-typeset/wordbreak"
)

type (
	// Option configures the Knuth-Plass line breaker.
	Option func(*options)

	options struct {
		tolerance float64
		badness   float64
		demerits  demeritsT
		formatterOptions
		looseness int // parameter q in the paper
	}

	formatterOptions struct {
		measurer    func([]rune) float64
		scaleFactor float64 // the scale factor is a multiplier to adapt measures to better fit the algorithm's settings
		space       sums    // space widths

		wordBreak          bool                  // enable breaking words (hyphenations, ...)
		renderHyphens      bool                  // enable the rendering of hyphens for hyphenated words
		hyphenPenalty      float64               // penalty to give to hyphenated words
		hardHyphenPenalty  float64               // penalty to give to explicitly hyphenated words
		punctuationPenalty float64               // penalty to give to punctuation marks
		hyphenator         wordbreaker.SplitFunc // word breaker for hyphenation
		punctuator         wordbreaker.SplitFunc // word breaker for punctuations signs (and more generally, all kind of "natural" separators)
		minHyphenate       int                   // minimum length of a token for hyphenation to apply
		glueStretch        float64
		glueShrink         float64
	}
)

// WithTolerance sets the threshold on an acceptable adjustment ratio.
//
// It corresponds to the parameter rho in the original paper.
//
// The default is 8.6.
func WithTolerance(tolerance float64) Option {
	return func(o *options) {
		o.tolerance = tolerance
	}
}

func WithScaleFactor(scale float64) Option {
	return func(o *options) {
		o.scaleFactor = scale
	}
}

// WithWordBreak enables single words (tokens) to be broken down.
//
// By default, this is enabled.
func WithWordBreak(enabled bool) Option {
	return func(o *options) {
		o.wordBreak = enabled
	}
}

// WithMeasurer sets a function to measure the width of a string.
//
// By default, all characters in a token are considered with width 1 (i.e. len(string)).
func WithMeasurer(measurer func([]rune) float64) Option {
	return func(o *options) {
		o.measurer = measurer
	}
}

// WithHyhenator specifies a SplitFunc operator to break down words.
//
// It implies WithWordBreak(true).
func WithHyphenator(hyphenator wordbreaker.SplitFunc) Option {
	return func(o *options) {
		o.wordBreak = true
		o.hyphenator = hyphenator
	}
}

// WithPunctuator specifies a SplitFunc operator to break down words that contain punctuation marks.
func WithPunctuator(punctuator wordbreaker.SplitFunc) Option {
	return func(o *options) {
		o.punctuator = punctuator
	}
}

// WithHyphenPenalty sets the penalty attributed to word breaks.
func WithHyphenPenalty(penalty float64) Option {
	return func(o *options) {
		o.hyphenPenalty = penalty
		o.demerits.flagged = penalty
	}
}

// WithRenderHypens enables the insertion of hyphens ("-") at the end of a line
// when rendering broken down words.
//
// This is enabled by default.
func WithRenderHyphens(enabled bool) Option {
	return func(o *options) {
		o.renderHyphens = enabled
	}
}

func WithLooseness(looseness int) Option {
	return func(o *options) {
		o.looseness = looseness
	}
}

func defaultOptions(opts []Option) *options {
	o := &options{
		tolerance:        8.6,
		badness:          100.00, // the badness constant from Knuth&Plass paper
		demerits:         defaultDemerits(),
		formatterOptions: defaultFormatterOptions(),
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

func defaultDemerits() demeritsT {
	return demeritsT{
		line:    10,  // line penalty (1.00 in Knuth & Plass)
		flagged: 100, // extra penalty applied to broken words. alpha parameter in the paper.
		fitness: 200, // gamma parameter in Knuth & Plass. Proposed with values 3000, 100
	}
}

func defaultFormatterOptions() formatterOptions {
	return formatterOptions{
		wordBreak:     true,
		renderHyphens: true,
		measurer:      func(in []rune) float64 { return float64(runes.Widths(in)) }, // TODO: strip ANSI escape seq
		scaleFactor:   1,                                                            // 3,
		space: sums{
			width:   1, // 3,
			stretch: 2, // 6,
			shrink:  3, // 9,
		},
		hyphenPenalty:      300, // penalty applied to breaks after a soft hyphen
		hardHyphenPenalty:  200, // penalty applied to breaks after an explicit hyphen
		punctuationPenalty: 400, // penalty applied to break before a punctuation mark
		minHyphenate:       4,   // minimum length of a word to be hyphenated
		glueStretch:        6,   // 12 -> 18,
		glueShrink:         0,   // ,
	}
}
