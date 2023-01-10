# go-typeset

Typesetting utilities for golang.

## Word breakers

* hyphenator: breaks a word across legit hyphenation breakpoints. Implements the classical algorithm from Frank M. Liang, with support for TeX hyphenation rule files.
* punctuator: breaks a word across punctuation marks, retaining separators
* tokenizer: breaks a text into space-separated tokens

### Maintenance

* hyphenator

The hyphenator uses embeded files, packed in memory at build time.
support for additional languages may be provided by copying the appropriate TeX hyphenation files.

Use `go generate ./...` to strip the files and make them tighter in memory.

### TODOs
* hyphenator
  * support a larger set of languages, possibly with a build tag guard to keep size as needed
  * investigate possibility for a faster trie. Trie access time is currently dominated by the performance of
    go map.

## Line breaker

Implements the classical paragraph breaking algorithm from D. Knuth  & M. Plass.

## Terminal utilities

Utilities to work with runes on a terminal.

* ansi: strips input from start/end ANSI escape sequences, retain the sequences.
* runes: calculates the width on display of runes, exposes a useful io.RuneReader, exposes a clone of `strings.FieldsFunc` suited to `[]rune` input.

### TODOs
* ansi: implement attributes renderer, so we may render ansi escape sequence across word or line breaks
* runes 

> **NOTE:** `runes.Width` heavily relied on `github.com/mattn/go-runewidth`, especially for East Asian support.
> However, I've disabled support for unicode graphemes, as too complex and significantly degrading performances.
>
> Perhaps in the future we may reintroduce a fasrer, allocation-free version of `github.com/rivo/uniseg` and iterate over unicode
> graphemes rather than runes.

## Performances

Most of the exposed utilities provide interfaces to work with slice of runes instead of strings.

> **Why?** 
> Because most of the processing on immutable strings can be achieved without extraneous memory allocations when using `[]rune` instead of `string`.
> This makes a big difference in performances.

I've added a few tests that collect memory and/or CPU profiles. Their output may be use to profile the desired functions
and find more optimizations.

## Credits and licenses

TODO(fred): credit authors, and mention their license terms

_Algorithms_

_Implementation_
* rune width 
  * github.com/mattn/go-runewidth
* hyphenation 
  * github.com/...
* line breaker

_Inspirators_
* paragraph breaking ...
