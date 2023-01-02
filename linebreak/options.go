package linebreak

import (
	"github.com/fredbi/go-typeset/wordbreak"
)

type (
	Option func(*options)

	options struct {
		tolerance float64
		badness   float64
		demerits  demeritsT
		formatterOptions
	}

	formatterOptions struct {
		measurer    func(string) float64 // TODO: func([]rune) int
		scaleFactor float64              // the scale factor is a multiplier to adapt measures to better fit the algorithm's settings
		space       sums                 // space widths

		wordBreak     bool                  // enable breaking words (hyphenations, ...)
		renderHyphens bool                  // enable the rendering of hyphens for hyphenated words
		hyphenPenalty float64               // penalty to give to hyphenated words
		hyphenator    wordbreaker.SplitFunc // word breaker for hyphenation TODO func([]rune) [][]rune ?
		minHyphenate  int                   // minimum length of a token for hyphenation to apply
		glueStretch   float64
		glueShrink    float64
	}
)

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

func WithWordBreak(enabled bool) Option {
	return func(o *options) {
		o.wordBreak = enabled
	}
}

func WithMeasurer(measurer func(string) float64) Option {
	return func(o *options) {
		o.measurer = measurer
	}
}

func WithHyphenator(hyphenator wordbreaker.SplitFunc) Option {
	return func(o *options) {
		o.wordBreak = true
		o.hyphenator = hyphenator
	}
}

func WithHyphenPenalty(penalty float64) Option {
	return func(o *options) {
		o.hyphenPenalty = penalty
		o.demerits.flagged = penalty
	}
}

func WithRenderHyphens(enabled bool) Option {
	return func(o *options) {
		o.renderHyphens = enabled
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
		line:    10,   // line penalty (1.00 in Knuth & Plass)
		flagged: 100,  // extra penalty applied to broken words
		fitness: 3000, // gamma parameter in Knuth & Plass
	}
}

func defaultFormatterOptions() formatterOptions {
	return formatterOptions{
		renderHyphens: true,
		measurer:      func(in string) float64 { return float64(len(in)) }, // TODO: use rune width and strip ANSI escape seq
		scaleFactor:   1,                                                   //3,
		space: sums{
			width:   1, // 3,
			stretch: 2, // 6,
			shrink:  3, // 9,
		},
		hyphenPenalty: 100, // penalty applied to hyphens
		hyphenator:    func(in string) []string { return []string{in} },
		minHyphenate:  4, // minimum length of a word to be hyphenated
		glueStretch:   6, // 12 -> 18,
		glueShrink:    0, // ,
	}
}
