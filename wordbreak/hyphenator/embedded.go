package hyphenator

import (
	"embed"
)

// Embeds in the build all *.tex pattern files in the "languages" folder.

//go:embed languages/*.tex
var texFS embed.FS
