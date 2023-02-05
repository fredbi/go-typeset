// Package runesio exposes a Reader and a Writer to work with slices of runes.
//
// SliceReader implements io.RuneReader and io.Seeker.
// The Writer knows how to write runes, like io.Write works with []byte.
package runesio
