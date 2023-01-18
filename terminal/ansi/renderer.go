package ansi

/*
var _ attributes.Renderer = &BoxRenderer{}

type (
	BoxRenderer struct {
		*attributes.Attribute
	}
)
*/

/*
// BoxRenderersFromToken extracts a collection of renderers from a token, with start/stop rendering sequences.
// It implements the attributes.Renderer interface.
//
// Since a single token may contain several sequences, there may be more renderers returned than tokens.
// TODO: deprecate
func BoxRenderersFromToken(token []rune) attributes.Renderers {
	attrs := make(attributes.Renderers, 0, 4)

	for {
		stripped := StripANSIFromRunes(token)
		attrs = append(attrs,
			NewBoxRenderer(stripped.Text, stripped.StartSequence, stripped.StopSequence),
		)
		if len(stripped.Remainder) == 0 {
			break
		}
	}

	return attrs
}
*/

func StripToken(token []rune) []StrippedToken {
	attrs := make([]StrippedToken, 0, 4)

	for {
		stripped := StripANSIFromRunes(token)
		attrs = append(attrs, stripped)
		if len(stripped.Remainder) == 0 {
			break
		}

		token = stripped.Remainder
	}

	return attrs
}

/*
// NewBoxRenderer builds a BoxRenderer from a raw token and start/stop sequences.
func NewBoxRenderer(stripped, start, stop []rune) *BoxRenderer {
	return &BoxRenderer{
		Attribute: attributes.New(stripped, start, stop),
	}
}
*/
