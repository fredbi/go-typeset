package runes

import "testing"

func BenchmarkWidth(b *testing.B) {
	// warm lookup cache
	_ = Width('A')

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		for _, testCase := range runewidthtests {
			_ = Width(testCase.in)
		}
	}
}

func BenchmarkWidths(b *testing.B) {
	// warm lookup cache
	_ = Width('A')
	str := []rune("string")

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		_ = Widths(str)
	}
}

func BenchmarkStringWidth(b *testing.B) {
	// warm lookup cache
	_ = Width('A')
	const str = "string"

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		_ = StringWidth(str)
	}
}

func BenchmarkBuildLookup(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		_ = buildLookupTable(optionsWithDefaults(nil))
	}
}
