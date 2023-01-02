// Package hyphenator is a word hyphenator.
//
// It provies a golang implementation of the word hyphenation algorithm
// proposed by Frank M. Liang (https://www.tug.org/docs/liang/), which is distributed
// with Tex.
//
// The hyphenator uses the same language-specific rules as those distributed with TeX.
//
// Rule files are available at https://www.tug.org/tex-hyphen.
//
// Authors may contribute to the rule sets there: https://github.com/hyphenation/tex-hyphen.
//
// Credits to Norbert Pillmayer, whose work largely inspired this package (https://github.com/npillmayer/gotype),
// in particular the patterns dictionary parser.
package hyphenator
