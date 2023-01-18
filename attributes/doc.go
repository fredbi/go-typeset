// Package attributes exposes an interface that knows
// how to deal with text presentation attributes (color, typeface, etc).
//
// A Renderer is a unitary block of text, with an (optional) presentation context (StartSequence, EndSequence).
//
// A sequence may typically be a terminal ANSI SGR escape sequence, but any tag language (e.g. HTML, LaTeX) works
// with a similar logic.
//
// The State is a chained list of Renderers to track the context of Start/Stop sequences.
package attributes
