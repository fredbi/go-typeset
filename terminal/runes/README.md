# runes

Provides the cell width of a unicode rune, with extra support for East-Asian rules and emojis.

This is used to determine how many "cells" (e.g. on a terminal display) a rune takes on a line.

Limitation: multi-runes unicode grapheme clusters are currently not supported.

## Maintainance

A built-in table of special rune properties is generated from the properties published at https://www.unicode.org.

To update mappings to newer unicode releases, please update the `build_tables.go` source and re-generate the tables with:

```
go generate ./...
```
