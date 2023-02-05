// Package runes provides utilities to work with runes.
//
// * determines the width on a terminal display of a rune.
// * FieldsFunc splits a slice of runes like strings.FieldsFunc
//
// The current version is based on mappings and properties defined by unicode v15.0.0.
//
// Rune width supports East-Asian runes, including when supporting special character sets with wide characters.
//
// NOTE: unicode grapheme clusters are not currently supported, meaning that the width of multi-runes graphemes such as "🏳️\u200d🌈"
// is not properly calculated.
// Use-cases with unicode graphemes might find this implemention useful: https://github.com/rivo/uniseg.
package runes
