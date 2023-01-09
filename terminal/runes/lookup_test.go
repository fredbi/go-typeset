package runes

import "testing"

func BenchmarkBuildLookup(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		_ = buildLookupTable(optionsWithDefaults(nil))
	}
}
