package hyphenator

import (
	"embed"
)

// Embeds in the build all *.tex pattern files in the "languages" folder.
//
// TODO(fredbi): we'll be able to extend language support with some extra build tags,
// and thus optionally embed more languages.

//go:embed languages/*.tex
var texFS embed.FS
