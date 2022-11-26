package fonts

import "embed"

//
//go:embed *.ttf
var DefaultEmbeddedFontsFS embed.FS

//go:embed vogue.ttf
var DefaultEmbeddedFontFS embed.FS
