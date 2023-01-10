package runes

// FieldsFunc splits a slice of runes like strings.FieldsFunc.
func FieldsFunc(in []rune, splitFunc func(rune) bool) [][]rune {
	tokens := make([][]rune, 0, len(in))
	start := -1

	for end, r := range in {
		if splitFunc(r) {
			if start >= 0 {
				tokens = append(tokens, in[start:end])
				start = ^start
			}

			continue
		}

		if start < 0 {
			start = end
		}
	}

	if start >= 0 {
		// last token
		tokens = append(tokens, in[start:])
	}

	return tokens[:len(tokens):len(tokens)]
}
