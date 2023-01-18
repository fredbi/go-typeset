package ansi

import "testing"

func BenchmarkStripANSIFromRunes(b *testing.B) {
	const input = startInput + wordInput + endInput
	str := []rune(input)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(0)

	for n := 0; n < b.N; n++ {
		_ = StripANSIFromRunes(str)
	}
}
