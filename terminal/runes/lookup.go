package runes

import (
	"sync"
)

type (
	interval struct {
		first rune
		last  rune
	}

	table []interval
)

var (
	nonprint = table{
		{0x0000, 0x001F}, {0x007F, 0x009F}, {0x00AD, 0x00AD},
		{0x070F, 0x070F}, {0x180B, 0x180E}, {0x200B, 0x200F},
		{0x2028, 0x202E}, {0x206A, 0x206F}, {0xD800, 0xDFFF},
		{0xFEFF, 0xFEFF}, {0xFFF9, 0xFFFB}, {0xFFFE, 0xFFFF},
	}

	buildLookupOnce sync.Once
	lookupTable     []byte
)

// initLookupTable allocates an in-memory lookup table of 278528 bytes for faster operation.
//
// Rune widths (<4) are packed on 2 bits.
func initLookupTable() {
	o := optionsWithDefaults(nil)

	lookupTable = buildLookupTable(o)
}

func buildLookupTable(o *options) []byte {
	const max = 0x110000

	lookup := make([]byte, max>>2)
	for i := range lookup {
		i32 := int32(i * 4)
		x0 := runeWidth(i32, o)   // rune with index % 4 = 0
		x1 := runeWidth(i32+1, o) // rune with index % 4 = 1
		x2 := runeWidth(i32+2, o) // rune with index % 4 = 2
		x3 := runeWidth(i32+3, o) // rune with index % 4 = 3

		// pack 4 widths (value<4) in a single byte
		lookup[i] = uint8(x0) |
			uint8(x1)<<2 |
			uint8(x2)<<4 |
			uint8(x3)<<6
	}

	return lookup
}

func inTables(r rune, ts ...table) bool {
	for _, t := range ts {
		if inTable(r, t) {
			return true
		}
	}

	return false
}

func inTable(r rune, t table) bool {
	if r < t[0].first {
		return false
	}

	bot := 0
	top := len(t) - 1

	for top >= bot {
		mid := (bot + top) >> 1

		switch {
		case t[mid].last < r:
			bot = mid + 1
		case t[mid].first > r:
			top = mid - 1
		default:
			return true
		}
	}

	return false
}
