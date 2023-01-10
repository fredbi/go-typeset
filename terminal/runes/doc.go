// Package runes provides utilities to work with runes.
//
// * determines the width on a terminal display of a rune.
// * exposes an io.RuneReader made out of a []rune slice.
// * splits a slice of runes like strings.FieldsFunc
//
// Rune width supports East-Asian runes, including when supporting special character sets with wide characters.
//
// NOTE: unicode grapheme clusters are not supported.
package runes
