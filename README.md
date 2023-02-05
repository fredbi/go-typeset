# go-typeset

Typesetting utilities for golang.

## Features

This package exposes tools to _break_ words and lines.

It is primarily intended to support text rendering on a terminal,
with monospace fonts and ANSI escape sequences.

## Project Status
The project is currently a WIP. ETA for a stabilized API: Feb. 2023

## Word breakers

Work breakers work on slices of runes.

* hyphenator: breaks a word across legit hyphenation breakpoints. Implements the classical algorithm from Frank M. Liang, with support for TeX hyphenation rule files.
* punctuator: breaks a word across punctuation marks, retaining separators
* tokenizer: breaks a text into space-separated tokens

> Notice that the hyphenator does not introduce the additional hyphens in the result.
> In order to distinguish such "soft-hyphens" from hyphens already in the input, we use `hyphenator.SplitWord()`

Code examples:

```golang
h := hyphenator.New(hyphenator.WithLangageTag(language.German))
parts := h.BreakWord([]rune("Ausnahme")

// [][]rune{
//   []rune("Aus"),
//   []rune("nah"),
//   []rune("me"),
// }
```

```golang
p := punctuator.New()
parts := p.BreakWord([]rune(`axy,b/cuv`)

// [][]rune{
//	[]rune("axy"), 
//  []rune(","), 
//  []rune("b"), 
//  []rune("/"), 
//  []rune("cuv"), 
// }
```

```golang
t := tokenizer.New()

tokens := t.BreakWord("In olden times when wishing still helped one")

// [][]rune{
//	[]rune("In"), 
//  []rune("olden"), 
//  []rune("times"), 
//  []rune("when"), 
//  []rune("whishing"), 
//  []rune("still"), 
//  []rune("helped"), 
//  []rune("one"), 
// }
```

## Line breaker

The line breaker implements the classical paragraph breaking algorithm from D. Knuth  & M. Plass,
to wrap words nicely under width and alignment constraints.

## Terminal utilities

Utilities to work with runes on a terminal.

* ansi: identifies and strips input from start/end ANSI escape sequences
* runes: calculates the width of runes on display
* runesio: reader & writer to manipulate slice of runes

## Rendering attributes on a terminal

## In-depth

### Dependencies
This repository comes with almost no non-stdlib dependencies, save for testing & profiling.

### Performances

Most of the exposed utilities provide interfaces to work with slice of runes instead of strings.

> **Why?** 
> Because most of the processing on immutable strings can be achieved without extraneous memory allocations when using `[]rune` instead of `string`.

I've added a few tests that collect memory and/or CPU profiles. Their output may be use to profile the desired functions
and find more optimizations.

### Maintenance

Some of the provided packages come with built-in external sources, that need to be updated from time to time. See:

* [runes widths](terminal/runes/README.md)
* [hyphenator](wordbreak/hyphenator/README.md)

### Limitations

* Fonts: at this moment, rendering utilities exposed by this package only support fixed-width terminal output.
* Runes: at this moment, rune width calculation does not support unicode grapheme clusters

> Perhaps in the future we may reintroduce a faster, allocation-free version of `github.com/rivo/uniseg` and iterate over unicode
> graphemes rather than runes.

### TODOs

A piece of software can never been considered as complete... Here are a few possible future directions.

* hyphenator
  * support a larger set of languages, possibly with a build tag guard to keep size as needed
  * investigate possibility for a faster trie. Trie access time is currently dominated by the performance of the inner go map.
* line breaker
  * future musings could extend the rendering to support PDF/HTML output, with font width measuring etc. Wow!
  * add the simpler greedy algorithm for comparison (e.g performance vs quality)

## Credits and licenses

_Algorithms_

* hyphenation: algorithm by Frank M. Liang
* line breaking: algorithm by Donald E. Knuth and Michael F. Plass

A copy of the original paper from Knuth&Plass is available [here](docs/breaking-paragraphs-into-lines-donald-e-knuth-and-michael-f-plass.pdf)

_Implementations_

* rune width calculation is largely inspired from https://github.com/mattn/go-runewidth, under the MIT license. Thank you @mattn (Yasuhiro Matsumoto).
* hyphenation borrows a few parts from the experimental work published by Norbert Pillmayer at https://github.com/npillmayer/gotype, under the BSD license.
* line breaker is an original implementation in golang. However, existing implementations in other languages proved much inspirational. Credits:
  * https://github.com/bramstein/typeset (javascript)
  * https://github.com/baskerville/paragraph-breaker (rust)
  * https://github.com/alex-panda/KnuthPlassLineBreak (python)
* the internal trie derives from the work from Dalton G. Hubble, originally shared at https://github.com/dghubble/trie under the MIT license.
